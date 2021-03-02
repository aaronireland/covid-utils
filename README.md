<p align="center">
<img src="https://raw.githubusercontent.com/ashleymcnamara/gophers/master/NERDY.png" width=175 />
<img src="https://d2q79iu7y748jz.cloudfront.net/s/_squarelogo/73ec9cb60ab678cef553b9d4b46744b4" width=150 />
</p>

<h1 align="center">Covid-19 Utilites</h1>

<p align="center"><i>CLI and API Clients for Covid-19-related data, written in Go.</i></p>

## Usage

Check for available vaccination appoints at Rite-Aid locations using the CLI

```shell
go build -o rite-aid-site ./cmd/vaccine-finder
./rite-aid-sites "123 Main St. Anytown, PA 17000" 11111,11112,11113

# restrict to only a certain state
./rite-aid-sites -state PA "123 Main St. Anytown, PA 17000" 12345
```

Add `-debug=true` to add verbose logging and pick a random site to text (for testing - this won't be a real available appointment)

```shell
./rite-aid-sites -debug=true -sms 2155550101 19002
```

Add SMS alerts using a [Twilio account](https://www.twilio.com/sms):

```shell
export TWILIO_ACCOUNT_SID="your-acccount-id"
export TWILIO_AUTH_TOKEN="your-auth-token"

./rite-aid-sites -sms "+15554206969" 12345,11234
```
