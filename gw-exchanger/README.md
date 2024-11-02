**gw-exchanger**

Сервис для получения курсов валют

structs:
-CurrencyRequest: [string] from_currency, [string] to_currency
-ExchangeRateResponse: [string] from_currency, [string] to_currenty, [float] rate
-ExchangeRatesResponse: [map<string, float>] rates
-Empty: nothing

sevices:
Основной сервис
service ExchangeService

Метод для получения курса всех валют
-GetExchangeRates(Empty) returns (ExchangeRatesResponse);
Метод для получения курса конкретной валюты
-GetExchangeRateForCurrency(CurrencyRequest) returns (ExchangeRateResponse);