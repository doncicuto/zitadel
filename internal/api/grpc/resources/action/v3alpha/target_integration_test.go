//go:build integration

package action_test

import (
	"context"
	"fmt"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	action "github.com/zitadel/zitadel/pkg/grpc/resources/action/v3alpha"
	resource_object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
)

func TestServer_CreateTarget(t *testing.T) {
	_, instanceID, _, isolatedIAMOwnerCTX := Tester.UseIsolatedInstance(t, IAMOwnerCTX, SystemCTX)
	ensureFeatureEnabled(t, isolatedIAMOwnerCTX)
	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.Target
		want    *resource_object.Details
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &action.Target{
				Name: fmt.Sprint(time.Now().UnixNano() + 1),
			},
			wantErr: true,
		},
		{
			name: "empty name",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "empty type",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name:       fmt.Sprint(time.Now().UnixNano() + 1),
				TargetType: nil,
			},
			wantErr: true,
		},
		{
			name: "empty webhook url",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name: fmt.Sprint(time.Now().UnixNano() + 1),
				TargetType: &action.Target_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{},
				},
			},
			wantErr: true,
		},
		{
			name: "empty request response url",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name: fmt.Sprint(time.Now().UnixNano() + 1),
				TargetType: &action.Target_RestCall{
					RestCall: &action.SetRESTCall{},
				},
			},
			wantErr: true,
		},
		{
			name: "empty timeout",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name:     fmt.Sprint(time.Now().UnixNano() + 1),
				Endpoint: "https://example.com",
				TargetType: &action.Target_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{},
				},
				Timeout: nil,
			},
			wantErr: true,
		},
		{
			name: "async, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name:     fmt.Sprint(time.Now().UnixNano() + 1),
				Endpoint: "https://example.com",
				TargetType: &action.Target_RestAsync{
					RestAsync: &action.SetRESTAsync{},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_INSTANCE,
					Id:   instanceID,
				},
			},
		},
		{
			name: "webhook, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name:     fmt.Sprint(time.Now().UnixNano() + 1),
				Endpoint: "https://example.com",
				TargetType: &action.Target_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{
						InterruptOnError: false,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_INSTANCE,
					Id:   instanceID,
				},
			},
		},
		{
			name: "webhook, interrupt on error, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name:     fmt.Sprint(time.Now().UnixNano() + 1),
				Endpoint: "https://example.com",
				TargetType: &action.Target_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{
						InterruptOnError: true,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_INSTANCE,
					Id:   instanceID,
				},
			},
		},
		{
			name: "call, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name:     fmt.Sprint(time.Now().UnixNano() + 1),
				Endpoint: "https://example.com",
				TargetType: &action.Target_RestCall{
					RestCall: &action.SetRESTCall{
						InterruptOnError: false,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_INSTANCE,
					Id:   instanceID,
				},
			},
		},

		{
			name: "call, interruptOnError, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name:     fmt.Sprint(time.Now().UnixNano() + 1),
				Endpoint: "https://example.com",
				TargetType: &action.Target_RestCall{
					RestCall: &action.SetRESTCall{
						InterruptOnError: true,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_INSTANCE,
					Id:   instanceID,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tester.Client.ActionV3Alpha.CreateTarget(tt.ctx, &action.CreateTargetRequest{Target: tt.req})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.want, got.Details)
		})
	}
}

func TestServer_PatchTarget(t *testing.T) {
	_, instanceID, _, isolatedIAMOwnerCTX := Tester.UseIsolatedInstance(t, IAMOwnerCTX, SystemCTX)
	ensureFeatureEnabled(t, isolatedIAMOwnerCTX)
	type args struct {
		ctx context.Context
		req *action.PatchTargetRequest
	}
	tests := []struct {
		name    string
		prepare func(request *action.PatchTargetRequest) error
		args    args
		want    *resource_object.Details
		wantErr bool
	}{
		{
			name: "missing permission",
			prepare: func(request *action.PatchTargetRequest) error {
				targetID := Tester.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetDetails().GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: Tester.WithAuthorization(context.Background(), integration.OrgOwner),
				req: &action.PatchTargetRequest{
					Target: &action.PatchTarget{
						Name: gu.Ptr(fmt.Sprint(time.Now().UnixNano() + 1)),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "not existing",
			prepare: func(request *action.PatchTargetRequest) error {
				request.Id = "notexisting"
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.PatchTargetRequest{
					Target: &action.PatchTarget{
						Name: gu.Ptr(fmt.Sprint(time.Now().UnixNano() + 1)),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "change name, ok",
			prepare: func(request *action.PatchTargetRequest) error {
				targetID := Tester.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetDetails().GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.PatchTargetRequest{
					Target: &action.PatchTarget{
						Name: gu.Ptr(fmt.Sprint(time.Now().UnixNano() + 1)),
					},
				},
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_INSTANCE,
					Id:   instanceID,
				},
			},
		},
		{
			name: "change type, ok",
			prepare: func(request *action.PatchTargetRequest) error {
				targetID := Tester.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetDetails().GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.PatchTargetRequest{
					Target: &action.PatchTarget{
						TargetType: &action.PatchTarget_RestCall{
							RestCall: &action.SetRESTCall{
								InterruptOnError: true,
							},
						},
					},
				},
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_INSTANCE,
					Id:   instanceID,
				},
			},
		},
		{
			name: "change url, ok",
			prepare: func(request *action.PatchTargetRequest) error {
				targetID := Tester.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetDetails().GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.PatchTargetRequest{
					Target: &action.PatchTarget{
						Endpoint: gu.Ptr("https://example.com/hooks/new"),
					},
				},
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_INSTANCE,
					Id:   instanceID,
				},
			},
		},
		{
			name: "change timeout, ok",
			prepare: func(request *action.PatchTargetRequest) error {
				targetID := Tester.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetDetails().GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.PatchTargetRequest{
					Target: &action.PatchTarget{
						Timeout: durationpb.New(20 * time.Second),
					},
				},
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_INSTANCE,
					Id:   instanceID,
				},
			},
		},
		{
			name: "change type async, ok",
			prepare: func(request *action.PatchTargetRequest) error {
				targetID := Tester.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeAsync, false).GetDetails().GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.PatchTargetRequest{
					Target: &action.PatchTarget{
						TargetType: &action.PatchTarget_RestAsync{
							RestAsync: &action.SetRESTAsync{},
						},
					},
				},
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_INSTANCE,
					Id:   instanceID,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.prepare(tt.args.req)
			require.NoError(t, err)
			// We want to have the same response no matter how often we call the function
			Tester.Client.ActionV3Alpha.PatchTarget(tt.args.ctx, tt.args.req)
			got, err := Tester.Client.ActionV3Alpha.PatchTarget(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.want, got.Details)
		})
	}
}

func TestServer_DeleteTarget(t *testing.T) {
	_, instanceID, _, isolatedIAMOwnerCTX := Tester.UseIsolatedInstance(t, IAMOwnerCTX, SystemCTX)
	ensureFeatureEnabled(t, isolatedIAMOwnerCTX)
	target := Tester.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false)
	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.DeleteTargetRequest
		want    *resource_object.Details
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &action.DeleteTargetRequest{
				Id: target.GetDetails().GetId(),
			},
			wantErr: true,
		},
		{
			name: "empty id",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.DeleteTargetRequest{
				Id: "",
			},
			wantErr: true,
		},
		{
			name: "delete target",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.DeleteTargetRequest{
				Id: target.GetDetails().GetId(),
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_INSTANCE,
					Id:   instanceID,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tester.Client.ActionV3Alpha.DeleteTarget(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.want, got.Details)
		})
	}
}
