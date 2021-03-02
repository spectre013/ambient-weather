package main

import (
	"fmt"

	// other imports
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func getClient(creds *Credentials) (*twitter.Client, error) {
	// Pass in your consumer key (API Key) and your Consumer Secret (API Secret)
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	// Pass in your Access Token and your Access Token Secret
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	// we can retrieve the user and verify if the credentials
	// we have used successfully allow us to log in!
	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}
	logger.Info("Logged ", user.Name, " into Twitter")
	return client, nil
}

func sendUpdate() {
	query := `select id,mac,date,baromabsin,baromrelin,battout,Batt1,Batt2,Batt3,Batt4,Batt5,Batt6,Batt7,Batt8,Batt9,Batt10,co2,dailyrainin,dewpoint,eventrainin,feelslike,
				hourlyrainin,hourlyrain,humidity,humidity1,humidity2,humidity3,humidity4,humidity5,humidity6,humidity7,humidity8,humidity9,humidity10,humidityin,lastrain,
				maxdailygust,relay1,relay2,relay3,relay4,relay5,relay6,relay7,relay8,relay9,relay10,monthlyrainin,solarradiation,tempf,temp1f,temp2f,temp3f,temp4f,temp5f,temp6f,temp7f,temp8f,temp9f,temp10f,
				tempinf,totalrainin,uv,weeklyrainin,winddir,windgustmph,windgustdir,windspeedmph,yearlyrainin,battlightning,lightningday,lightninghour,lightningtime,lightningdistance 
				from records order by date desc limit 0,1`
	rec := getRecord(query)
	t := buildMessage(rec)
	_, _, err := client.Statuses.Update(t, nil)
	if err != nil {
		logger.Println(err)
	}
}

func buildMessage(rec Record) string {
	t := ""
	if rec.Tempf < 50 {
		t = fmt.Sprintf("Temp: %.2f°F Sustained winds at %.2fmph, gusts to: %.2fmph, Feels like %.2f #COwx #KCOCOLOR663", rec.Tempf, rec.Windspeedmph, rec.Windgustmph, rec.Feelslike)
	} else {
		t = fmt.Sprintf("Temp: %.2f°F Sustained winds at %.2fmph, gusts to: %.2fmph, Rain: %.2fin Lightning: %d Today #COwx #KCOCOLOR663", rec.Tempf, rec.Windspeedmph, rec.Windgustmph, rec.Dailyrainin, rec.Lightningday)
	}
	return t
}
