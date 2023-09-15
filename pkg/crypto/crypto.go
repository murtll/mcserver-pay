package crypto

import (
	"crypto/md5"
	"fmt"
)

func CheckSignFk(args ...interface{}) (bool, error) {
	if len(args) != 5 {
		return false, fmt.Errorf("expected 5 args to check fk sign, got %d", len(args))
	}
	sign := args[0]
	merchantId := args[1]
	amount := args[2]
	merchantSecret := args[3]
	merchantOrderId := args[4]
	expectedSign := string(md5.New().Sum([]byte(fmt.Sprintf("%v:%v:%v:%v", merchantId, amount, merchantSecret, merchantOrderId))))
	return sign == expectedSign, nil
}
