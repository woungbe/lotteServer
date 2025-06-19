package load

/*
	// 초기 데이터 가져오기 !!




*/

import (
	"log"

	lotto "github.com/woungbe/lotto-fetcher/pkg/lotto"
)

// 회사 번호로 데이터 가져오기 !!
func GetLottoRound(number int) (*lotto.LottoResult, error) {
	res, err := lotto.FetchLotto(number)

	if err != nil {
		log.Println(err)
	}
	return res,nil
}