package eventconsume

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/hublabs/common/auth"

	"github.com/go-xorm/xorm"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/ctxdb"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/pangpanglabs/goutils/kafka"
)

type FilterFunc func(ctx context.Context) error
type Filter func(next FilterFunc) FilterFunc

func Handle(ctx context.Context, filters []Filter, f FilterFunc) error {
	for i := range filters {
		f = filters[len(filters)-1-i](f)
	}
	return f(ctx)
}

func ContextDB(service string, xormEngine *xorm.Engine, kafkaConfig kafka.Config) Filter {
	return ContextDBWithName(service, echomiddleware.ContextDBName, xormEngine, kafkaConfig)
}

func ContextTxDB(service string, xormEngine *xorm.Engine, kafkaConfig kafka.Config) Filter {
	return ContextTxDBWithName(service, echomiddleware.ContextDBName, xormEngine, kafkaConfig)
}

func ContextDBWithName(service string, contexDBName echomiddleware.ContextDBType, xormEngine *xorm.Engine, kafkaConfig kafka.Config) Filter {
	ctxdb := ctxdb.New(xormEngine, service, kafkaConfig)

	return func(next FilterFunc) FilterFunc {
		return func(ctx context.Context) error {
			session := ctxdb.NewSession(ctx)
			defer session.Close()

			ctx = context.WithValue(ctx, contexDBName, session)

			return next(ctx)
		}
	}
}

func ContextTxDBWithName(service string, contexDBName echomiddleware.ContextDBType, xormEngine *xorm.Engine, kafkaConfig kafka.Config) Filter {
	ctxdb := ctxdb.New(xormEngine, service, kafkaConfig)

	return func(next FilterFunc) FilterFunc {
		return func(ctx context.Context) error {
			session := ctxdb.NewSession(ctx)
			defer session.Close()

			if err := session.Begin(); err != nil {
				log.Println("Tx Begin Error.", err)
			}

			ctx = context.WithValue(ctx, contexDBName, session)

			if err := next(ctx); err != nil {
				if err := session.Rollback(); err != nil {
					log.Println("Tx Rollback Error.", err)
				}
				return err
			}

			if err := session.Commit(); err != nil {
				log.Println("Tx Commit Error.", err)
				return err
			}

			return nil
		}
	}
}

func Recover() Filter {
	stackSize := 4 << 10 // 4 KB
	disableStackAll := false
	disablePrintStack := false

	return func(next FilterFunc) FilterFunc {
		return func(ctx context.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, stackSize)
					length := runtime.Stack(stack, !disableStackAll)
					if !disablePrintStack {
						log.Printf("[PANIC RECOVER] %v %s\n", err, stack[:length])
					}
					behaviorlog.FromCtx(ctx).WithError(err)
				}
			}()

			return next(ctx)
		}
	}
}

func UserClaimMiddleware() Filter {
	return func(next FilterFunc) FilterFunc {

		return func(ctx context.Context) error {

			userClaim := auth.UserClaim{}.FromCtx(ctx)
			if userClaim.Issuer != "" {
				return next(ctx)
			}

			return next(context.WithValue(ctx, "userClaim", newUserClaimFromToken(behaviorlog.FromCtx(ctx).AuthToken)))
		}
	}
}
func newUserClaimFromToken(token string) (userClaim auth.UserClaim) {
	userClaim, _ = auth.UserClaim{}.FromToken(token)
	return
}
