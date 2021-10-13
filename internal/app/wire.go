// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package app

import (
	"github.com/google/wire"

	"email-send-manager/internal/app/api"
	"email-send-manager/internal/app/dao"
	"email-send-manager/internal/app/router"
	"email-send-manager/internal/app/service"
)

// BuildInjector 生成注入器
func BuildInjector() (*Injector, func(), error) {
	wire.Build(
		InitGormDB,
		dao.RepoSet,
		InitGinEngine,
		service.ServiceSet,
		api.APISet,
		router.RouterSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
