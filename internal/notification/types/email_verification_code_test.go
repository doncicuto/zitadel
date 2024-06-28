package types

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestNotify_SendEmailVerificationCode(t *testing.T) {
	type args struct {
		user          *query.NotifyUser
		origin        string
		code          string
		urlTmpl       string
		authRequestID string
		loginPolicy   *query.LoginPolicy
	}
	tests := []struct {
		name    string
		args    args
		want    *notifyResult
		wantErr error
	}{
		{
			name: "default URL",
			args: args{
				user: &query.NotifyUser{
					ID:            "user1",
					ResourceOwner: "org1",
				},
				origin:        "https://example.com",
				code:          "123",
				urlTmpl:       "",
				authRequestID: "authRequestID",
				loginPolicy: &query.LoginPolicy{
					AllowUsernamePassword:      true,
					AllowRegister:              true,
					AllowExternalIDPs:          true,
					ForceMFA:                   true,
					ForceMFALocalOnly:          true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					HidePasswordReset:          true,
					IgnoreUnknownUsernames:     true,
					AllowDomainDiscovery:       true,
					DisableLoginWithEmail:      true,
					DisableLoginWithPhone:      true,
					DefaultRedirectURI:         "",
					PasswordCheckLifetime:      database.Duration(time.Hour),
					ExternalLoginCheckLifetime: database.Duration(time.Minute),
					MFAInitSkipLifetime:        database.Duration(time.Millisecond),
					SecondFactorCheckLifetime:  database.Duration(time.Microsecond),
					MultiFactorCheckLifetime:   database.Duration(time.Nanosecond),
					SecondFactors: []domain.SecondFactorType{
						domain.SecondFactorTypeTOTP,
						domain.SecondFactorTypeU2F,
						domain.SecondFactorTypeOTPEmail,
						domain.SecondFactorTypeOTPSMS,
					},
					MultiFactors: []domain.MultiFactorType{
						domain.MultiFactorTypeU2FWithPIN,
					},
					IsDefault: true,
					UseDefaultRedirectUriForNotificationLinks: false,
				},
			},
			want: &notifyResult{
				url:                                "https://example.com/ui/login/mail/verification?authRequestID=authRequestID&code=123&orgID=org1&userID=user1",
				args:                               map[string]interface{}{"Code": "123"},
				messageType:                        domain.VerifyEmailMessageType,
				allowUnverifiedNotificationChannel: true,
			},
		},
		{
			name: "template error",
			args: args{
				user: &query.NotifyUser{
					ID:            "user1",
					ResourceOwner: "org1",
				},
				origin:        "https://example.com",
				code:          "123",
				urlTmpl:       "{{",
				authRequestID: "authRequestID",
				loginPolicy: &query.LoginPolicy{
					AllowUsernamePassword:      true,
					AllowRegister:              true,
					AllowExternalIDPs:          true,
					ForceMFA:                   true,
					ForceMFALocalOnly:          true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					HidePasswordReset:          true,
					IgnoreUnknownUsernames:     true,
					AllowDomainDiscovery:       true,
					DisableLoginWithEmail:      true,
					DisableLoginWithPhone:      true,
					DefaultRedirectURI:         "",
					PasswordCheckLifetime:      database.Duration(time.Hour),
					ExternalLoginCheckLifetime: database.Duration(time.Minute),
					MFAInitSkipLifetime:        database.Duration(time.Millisecond),
					SecondFactorCheckLifetime:  database.Duration(time.Microsecond),
					MultiFactorCheckLifetime:   database.Duration(time.Nanosecond),
					SecondFactors: []domain.SecondFactorType{
						domain.SecondFactorTypeTOTP,
						domain.SecondFactorTypeU2F,
						domain.SecondFactorTypeOTPEmail,
						domain.SecondFactorTypeOTPSMS,
					},
					MultiFactors: []domain.MultiFactorType{
						domain.MultiFactorTypeU2FWithPIN,
					},
					IsDefault: true,
					UseDefaultRedirectUriForNotificationLinks: false,
				},
			},
			want:    &notifyResult{},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOMAIN-oGh5e", "Errors.User.InvalidURLTemplate"),
		},
		{
			name: "template success",
			args: args{
				user: &query.NotifyUser{
					ID:            "user1",
					ResourceOwner: "org1",
				},
				origin:        "https://example.com",
				code:          "123",
				urlTmpl:       "https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
				authRequestID: "authRequestID",
				loginPolicy: &query.LoginPolicy{
					AllowUsernamePassword:      true,
					AllowRegister:              true,
					AllowExternalIDPs:          true,
					ForceMFA:                   true,
					ForceMFALocalOnly:          true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					HidePasswordReset:          true,
					IgnoreUnknownUsernames:     true,
					AllowDomainDiscovery:       true,
					DisableLoginWithEmail:      true,
					DisableLoginWithPhone:      true,
					DefaultRedirectURI:         "",
					PasswordCheckLifetime:      database.Duration(time.Hour),
					ExternalLoginCheckLifetime: database.Duration(time.Minute),
					MFAInitSkipLifetime:        database.Duration(time.Millisecond),
					SecondFactorCheckLifetime:  database.Duration(time.Microsecond),
					MultiFactorCheckLifetime:   database.Duration(time.Nanosecond),
					SecondFactors: []domain.SecondFactorType{
						domain.SecondFactorTypeTOTP,
						domain.SecondFactorTypeU2F,
						domain.SecondFactorTypeOTPEmail,
						domain.SecondFactorTypeOTPSMS,
					},
					MultiFactors: []domain.MultiFactorType{
						domain.MultiFactorTypeU2FWithPIN,
					},
					IsDefault: true,
					UseDefaultRedirectUriForNotificationLinks: false,
				},
			},
			want: &notifyResult{
				url:                                "https://example.com/email/verify?userID=user1&code=123&orgID=org1",
				args:                               map[string]interface{}{"Code": "123"},
				messageType:                        domain.VerifyEmailMessageType,
				allowUnverifiedNotificationChannel: true,
			},
		},
		{
			name: "use default uri for url link success",
			args: args{
				user: &query.NotifyUser{
					ID:            "user1",
					ResourceOwner: "org1",
				},
				origin:        "https://example.com",
				code:          "123",
				urlTmpl:       "https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
				authRequestID: "authRequestID",
				loginPolicy: &query.LoginPolicy{
					AllowUsernamePassword:      true,
					AllowRegister:              true,
					AllowExternalIDPs:          true,
					ForceMFA:                   true,
					ForceMFALocalOnly:          true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					HidePasswordReset:          true,
					IgnoreUnknownUsernames:     true,
					AllowDomainDiscovery:       true,
					DisableLoginWithEmail:      true,
					DisableLoginWithPhone:      true,
					DefaultRedirectURI:         "https://example.com",
					PasswordCheckLifetime:      database.Duration(time.Hour),
					ExternalLoginCheckLifetime: database.Duration(time.Minute),
					MFAInitSkipLifetime:        database.Duration(time.Millisecond),
					SecondFactorCheckLifetime:  database.Duration(time.Microsecond),
					MultiFactorCheckLifetime:   database.Duration(time.Nanosecond),
					SecondFactors: []domain.SecondFactorType{
						domain.SecondFactorTypeTOTP,
						domain.SecondFactorTypeU2F,
						domain.SecondFactorTypeOTPEmail,
						domain.SecondFactorTypeOTPSMS,
					},
					MultiFactors: []domain.MultiFactorType{
						domain.MultiFactorTypeU2FWithPIN,
					},
					IsDefault: true,
					UseDefaultRedirectUriForNotificationLinks: true,
				},
			},
			want: &notifyResult{
				url:                                "https://example.com",
				args:                               map[string]interface{}{"Code": "123"},
				messageType:                        domain.VerifyEmailMessageType,
				allowUnverifiedNotificationChannel: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, notify := mockNotify()
			err := notify.SendEmailVerificationCode(http_utils.WithComposedOrigin(context.Background(), tt.args.origin), tt.args.user, tt.args.code, tt.args.urlTmpl, tt.args.authRequestID, tt.args.loginPolicy)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
