package tenant

import (
	"wchatv1/proto/tenant"
	"context"
)


var _ tenant.TenantServiceHandler = &Service{}

type Service struct {}

func (s *Service) Create(ctx context.Context, req *tenant.CreateReq, rsp *tenant.CreateRsp) error {
	panic("implement me")
}

func (s *Service) Check(ctx context.Context, req *tenant.CheckReq, rsp *tenant.CheckRsp) error {
	panic("implement me")
}

//func (s *Service) Create(ctx context.Context, req *tenant.CreateReq) (*tenant.CreateRsp, error) {
//	rsp := &tenant.CreateRsp{}
//	tenant := &models.Tenants{
//		TenantCode: &req.TenantCode,
//		TenantKey:  &req.TenantKey,
//	}
//	res := utils.DB.Model(tenant).Create(tenant)
//	if res.Error != nil {
//		rsp.Code = 1
//		rsp.Message = res.Error.Error()
//		return rsp, nil
//	}
//
//	return rsp, nil
//}
//
//func (s *Service) Check(ctx context.Context, req *tenant.CheckReq) (*tenant.CheckRsp, error) {
//	rsp := &tenant.CheckRsp{}
//	tenant := &models.Tenants{
//		TenantCode: &req.TenantCode,
//		TenantKey:  &req.TenantKey,
//	}
//	res := utils.DB.Model(tenant).First(tenant)
//	if res.Error != nil {
//		if res.RecordNotFound() {
//			rsp.Code = 1
//			rsp.Message = fmt.Sprintf("租户信息错误")
//			return rsp, nil
//		} else {
//			return nil, res.Error
//		}
//	}
//
//	return rsp, nil
//}
