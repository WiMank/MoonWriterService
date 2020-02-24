package repository

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

func Test() {
	ctx := context.Background()
	str1 := "E:\\Dev\\GoDev\\MoonWriterService\\credentials.json"
	androidpublisherService, err := androidpublisher.NewService(ctx, option.WithCredentialsFile(str1))
	if err != nil {
		log.Infof("androidpublisherService ERR: ", err)
	}

	r, errr := androidpublisherService.Purchases.Products.Get("com.mwriter.moonwriter",
		"lemonade35",
		"bgeddnkemanoelbedjokoocc.AO-J1Owsddv6PUW4Ct4TWhiPqs0HgiL0wIBIjZgoWWaKPF9_nbti33qJQcSMzZFcBhrM-Lu7WJORZZr4m3C6iZ_wLuGyFLFp6UDTWR9syP27IAGq0lNo5NHtgNgtIXXTRISSW2g275ig").Do()
	if errr != nil {
		log.Infof("errr: ", errr)
	}
	log.Infof("REsult: ", r)

}
