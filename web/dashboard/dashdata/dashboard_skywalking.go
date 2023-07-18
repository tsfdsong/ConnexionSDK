package dashdata

import (
	"context"
	"errors"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/distmiddleware/skywalking"

	"github.com/SkyAPM/go2sky"
)

type RiskControlSubSpan struct {
	ctx     context.Context
	span    go2sky.Span
	TraceID string
}

func (r *RiskControlSubSpan) Handle(ctx context.Context, seq int) error {
	span, subctx, err := skywalking.CreateLocalSpan(ctx)
	if err != nil {
		return errors.New("WithdrawGameERC20Token erc20 withdraw CreateLocalSpan failed")
	}
	span.SetOperationName("FTWithdraw")
	r.ctx = subctx
	r.span = span

	reportSpan, ok := span.(go2sky.ReportedSpan)
	if !ok {
		return errors.New("WithdrawGameERC20Token erc20 withdraw span type is wrong")
	}

	r.TraceID = reportSpan.Context().TraceID
	go2sky.PutCorrelation(r.ctx, "trace_id", r.TraceID)
	go2sky.PutCorrelation(r.ctx, "type", fmt.Sprintf("%d", skywalking.FTRiskAutoCode))
	err = skywalking.SkyPutCorrelation(config.GetSkyWalkingConfig().KeyWithdraw, seq, config.GetSkyWalkingConfig().ValueWithdraw, r.ctx)
	if err != nil {
		return errors.New("WithdrawGameERC20Token erc20 withdraw SkyPutCorrelation failed")
	}
	return nil
}

func (r *RiskControlSubSpan) GetContext() context.Context {
	return r.ctx
}

func (r *RiskControlSubSpan) GetSpan() go2sky.Span {
	return r.span
}

func (r *RiskControlSubSpan) End() {
	r.span.End()
}

type GameServerSubSpan struct {
	ctx  context.Context
	span go2sky.Span
}

func (r *GameServerSubSpan) Handle(ctx context.Context, seq int) error {
	//create game server withdraw sub span
	gamespan, gamectx, err := skywalking.CreateLocalSpan(ctx)
	if err != nil {
		return errors.New("FT withdraw GamePrewithdraw CreateLocalSpan game server failed")
	}
	gamespan.SetOperationName("FTGameWithdraw")
	r.ctx = gamectx
	r.span = gamespan

	err = skywalking.SkyPutCorrelation(config.GetSkyWalkingConfig().KeyGameWithdraw, seq, config.GetSkyWalkingConfig().ValueGameWithdraw, r.ctx)
	if err != nil {
		return errors.New("FT withdraw GamePrewithdraw SkyPutCorrelation game server failed")
	}
	return nil
}

func (r *GameServerSubSpan) GetContext() context.Context {
	return r.ctx
}

func (r *GameServerSubSpan) GetSpan() go2sky.Span {
	return r.span
}

func (r *GameServerSubSpan) End() {
	r.span.End()
}
