package router

import (
	"context"

	"github.com/ao9911/bluebell-new/conf"
)

func Init(c *conf.Config) {
	startHttp(c)
}

func Stop(ctx context.Context) {
	stopHttp(ctx)
}
