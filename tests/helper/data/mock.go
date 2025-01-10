package data

import (
	"testing"
	"time"

	"github.com/banggok/boillerplate_architecture/internal/config/db"
	valueobject "github.com/banggok/boillerplate_architecture/internal/data/entity/value_object"
	"github.com/banggok/boillerplate_architecture/internal/data/model"
	"github.com/banggok/boillerplate_architecture/internal/pkg/password"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var TenantData = model.Tenant{
	Name:         "Doe's Bakery",
	Address:      "123 Main Street, Springfield, USA",
	Email:        "business@example.com",
	Phone:        "+0987654321",
	Timezone:     "America/New_York",
	OpeningHours: "08:00",
	ClosingHours: "20:00",
}

var PlainPassword = "password123"

var HashedPassword string

var AccountData = model.Account{
	Name:  "John Doe",
	Email: "rtriasmono@gmail.com",
	Phone: "+1234567890",
}

var Token string
var Now time.Time

func init() {
	token, _ := password.GeneratePassword(16)
	Token = *token
	Now = time.Now().UTC()
	AccountVerificationData.Token = &Token
	AccountVerificationData.ExpiresAt = Now.Add(24 * time.Hour)

	hashed, _ := password.HashPassword(PlainPassword)
	HashedPassword = *hashed
	AccountData.Password = HashedPassword
}

var AccountVerificationData = model.AccountVerification{
	Type: valueobject.EMAIL_VERIFICATION.String(),
}

func PrepareData[M any](t *testing.T, mysqlCfg *db.DBConnection, data M) {

	if err := mysqlCfg.Master.Session(&gorm.Session{FullSaveAssociations: true}).Save(&data).
		Error; err != nil {
		require.NoError(t, err)
	}

}
