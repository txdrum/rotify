package timer
 
import ( "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/line/line-bot-sdk-go/linebot"
    "github.com/mmcdole/gofeed"
)

type PubSubMessage struct {
    Data []byte `json:"data"`
}

func main() {
    dummy := PubSubMessage {}
    PubsubMain(nil, dummy)
}

func PubsubMain(ctx context.Context, m PubSubMessage) error {
    log.Println("PubsubMain begin")

    fp := gofeed.NewParser()
 
    feed, _ := fp.ParseURL("https://ragnarokonline.gungho.jp/index.rdf")

    items := feed.Items
    now := time.Now()
    format := "2006-01-02 15:04:05 MST"
    format2 := "2006-01-02T15:04:05-07:00"
    prevTime := now.Add(-1 * time.Hour)
    jst := time.FixedZone("Asia/Tokyo", 9 * 60 * 60)
    jstPrev := prevTime.In(jst)
    log.Printf("now: %s, prev: %s, jstPrev: %s", now.Format(format), prevTime.Format(format),
        jstPrev.Format(format))
    count := 0
    text := ""

    for _, item := range items {
        pubStr := item.Published
        pubTime, _ := time.Parse(format2, pubStr)

        log.Printf("Published: %s", pubTime.Format(format))

        if (jstPrev.Before(pubTime)) {
            log.Printf("Title: %s, Link: %s, Updated: %s\n",
                item.Title, item.Link, item.Updated)

            if (text == "") {
                text = item.Title + "\n" + item.Link + "\n"
            }

            count = count + 1
        }
    }

    log.Printf("count: %d", count)

    if (count == 0) {
        return nil
    } else if (count > 1) {
        text = text + fmt.Sprintf("他、%d件の更新があります。", (count - 1))
    }

    secret := os.Getenv("LINE_SECRET")
    token := os.Getenv("LINE_ACCESS_TOKEN")
    client, err := linebot.New(secret, token)

    if err != nil {
        log.Printf("error: %v", err)
        return err
    }

    log.Printf("text:%s", text)

    message := linebot.NewTextMessage(text)
    res, err := client.BroadcastMessage(message).Do()
    log.Printf("response: %v", res)

    if (err != nil) {
        log.Printf("error: %v", err)
        return err
    }

    log.Println("PubsubMain end")
    return nil
}
