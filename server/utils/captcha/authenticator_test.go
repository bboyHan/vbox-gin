package captcha

import (
	"testing"
)

func TestAuthQrCode(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name    string
		args    args
		url     string
		wantErr bool
	}{
		{
			name:    "Gen",
			args:    args{"zhangsan"},
			wantErr: false,
		},
		{
			name:    "Valid",
			args:    args{"882197"},
			wantErr: false,
		},
		{
			name:    "parse",
			args:    args{"882197"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Gen" {
				url, err := AuthQrCode(tt.args.d)
				t.Logf("GET() got = %v", url)
				if err != nil {
					t.Errorf("GET() got = %v", err)
				}
			}
			if tt.name == "Valid" {
				//otpauth://totp/VBOX:zhang%20san?algorithm=SHA1&digits=6&issuer=VBOX&period=30&secret=PWLCKLXIDWNJQFBI4ZNVVMUI2GHKGFEQ
				t.Logf("GET() got = %v", tt.args.d)
				ValidateCode("SUGNDGOULWYCIYZ7K3DIQWNDZOOKHWFE", tt.args.d)
			}
			if tt.name == "parse" {
				//otpauth://totp/VBOX:zhang%20san?algorithm=SHA1&digits=6&issuer=VBOX&period=30&secret=PWLCKLXIDWNJQFBI4ZNVVMUI2GHKGFEQ
				t.Logf("GET() got = %v", tt.args.d)
				ParseRemoteQrCodeImage("https://mass.alipay.com/wsdk/img?fileid=A*tMYGT5iWgA8AAAAAAAAAAAAACTInAQ&bz=payment_codec&zoom=original")
			}
		})
	}
}
