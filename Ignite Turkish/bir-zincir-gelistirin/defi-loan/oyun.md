# Oyun

Alice'in hesabına `10000foocoin` ekleyin. Bu tokenlar kredi teminatı olarak kullanılacaktır.

config.yml

```
version: 1
accounts:
  - name: alice
    coins:
      - 20000token
      - 10000foocoin
      - 200000000stake
  - name: bob
    coins:
      - 10000token
      - 100000000stake
client:
  openapi:
    path: docs/static/openapi.yml
faucet:
  name: bob
  coins:
    - 5token
    - 100000stake
validators:
  - name: alice
    bonded: 100000000stake
```

Bir blok zinciri node'u başlatın:

```
ignite chain serve
```

### Bir kredinin geri ödenmesi

Alice'in hesabından ücret olarak `100token` ve teminat olarak `1000foocoin` ile `1000token` kredi talep edin. Son tarih `500` blok olarak belirlenmiştir:

```
loand tx loan request-loan 1000token 100token 1000foocoin 500 --from alice
```

```
Loan:
- amount: 1000token
  borrower: cosmos153dk8qh56v4yg6e4uzrvvqjueu6d36fptlr2kw
  collateral: 1000foocoin
  deadline: "500"
  fee: 100token
  id: "0"
  lender: ""
  state: requested
```

Lütfen terminal pencerenizde görüntülenen adreslerin (`borrower` alanındakiler gibi) bu eğitimde verilenlerle eşleşmeyeceğini unutmayın. Bunun nedeni, bir hesabın `config.yml` dosyasında belirtilmiş bir anımsatıcısı olmadığı sürece Ignite'ın her zincir başlatıldığında yeni hesaplar oluşturmasıdır.

Approve the loan from Bob's account:

```
loand tx loan approve-loan 0 --from bob
```

```
loand q loan list-loan
```

`lender` alanı Bob'un adresi olarak güncellendi ve `state` `approved` olarak güncellendi:

```
Loan:
- amount: 1000token
  borrower: cosmos153dk8qh56v4yg6e4uzrvvqjueu6d36fptlr2kw
  collateral: 1000foocoin
  deadline: "500"
  fee: 100token
  id: "0"
  lender: cosmos1qfzpxfhsu2qfy2exkukuanrkzrrexh9yeg2pr4
  state: approved
```

```
loand q bank balances $(loand keys show alice -a)
```

`Foocoin` bakiyesi `9000` olarak güncellendi, çünkü `1000foocoin` modül hesabına teminat olarak aktarıldı. `Token` bakiyesi `21000` olarak güncellendi, çünkü `1000token` Bob'un hesabından Alice'in hesabına kredi olarak aktarıldı:

```
balances:
- amount: "9000"
  denom: foocoin
- amount: "100000000"
  denom: stake
- amount: "21000"
  denom: token
```

```
loand q bank balances $(loand keys show bob -a)  
```

`Token` bakiyesi `9000` olarak güncellenmiştir, çünkü `1000token` Bob'un hesabından Alice'in hesabına kredi olarak aktarılmıştır:

```
balances:
- amount: "100000000"
  denom: stake
- amount: "9000"
  denom: token
```

Alice'in hesabından krediyi geri ödeyin:

```
loand tx loan repay-loan 0 --from alice
```

```
loand q loan list-loan
```

`state` alanı `repayed` olarak güncellenmiştir:

```
Loan:
- amount: 1000token
  borrower: cosmos153dk8qh56v4yg6e4uzrvvqjueu6d36fptlr2kw
  collateral: 1000foocoin
  deadline: "500"
  fee: 100token
  id: "0"
  lender: cosmos1qfzpxfhsu2qfy2exkukuanrkzrrexh9yeg2pr4
  state: repayed
```

```
loand q bank balances $(loand keys show alice -a)
```

`Foocoin` bakiyesi `10000` olarak güncellendi, çünkü `1000foocoin` modül hesabından Alice'in hesabına aktarıldı. `Token` bakiyesi `19900` olarak güncellendi, çünkü `1000token` Alice'in hesabından Bob'un hesabına geri ödeme olarak aktarıldı ve `100token` Alice'in hesabından Bob'un hesabına ücret olarak aktarıldı:

```
balances:
- amount: "10000"
  denom: foocoin
- amount: "100000000"
  denom: stake
- amount: "19900"
  denom: token
```

```
loand q bank balances $(loand keys show bob -a)  
```

`Token` bakiyesi `10100` olarak güncellendi, çünkü `1000token` Alice'in hesabından Bob'un hesabına geri ödeme olarak aktarıldı ve `100token` Alice'in hesabından Bob'un hesabına ücret olarak aktarıldı:

```
balances:
- amount: "100000000"
  denom: stake
- amount: "10100"
  denom: token
```

### Bir kredinin tasfiyesi

Alice'in hesabından ücret olarak `100token` ve teminat olarak `1000foocoin` ile `1000token` kredi talep edin. Son tarih `20` blok olarak belirlenmiştir. Son tarih çok küçük bir değere ayarlanmıştır, böylece kredi bir sonraki adımda hızlı bir şekilde tasfiye edilebilir:

```
loand tx loan request-loan 1000token 100token 1000foocoin 20 --from alice
```

Listeye bir kredi eklenmiştir:

```
Loan:
- amount: 1000token
  borrower: cosmos153dk8qh56v4yg6e4uzrvvqjueu6d36fptlr2kw
  collateral: 1000foocoin
  deadline: "500"
  fee: 100token
  id: "0"
  lender: cosmos1qfzpxfhsu2qfy2exkukuanrkzrrexh9yeg2pr4
  state: repayed
- amount: 1000token
  borrower: cosmos153dk8qh56v4yg6e4uzrvvqjueu6d36fptlr2kw
  collateral: 1000foocoin
  deadline: "20"
  fee: 100token
  id: "1"
  lender: ""
  state: requested
```

Bob'un hesabından alınan krediyi onaylayın:

```
loand tx loan approve-loan 1 --from bob
```

Bob'un hesabındaki krediyi tasfiye edin:

```
loand tx loan liquidate-loan 1 --from bob
```

```
loand q loan list-loan
```

Kredi tasfiye edilmiştir:

```
Loan:
- amount: 1000token
  borrower: cosmos153dk8qh56v4yg6e4uzrvvqjueu6d36fptlr2kw
  collateral: 1000foocoin
  deadline: "500"
  fee: 100token
  id: "0"
  lender: cosmos1qfzpxfhsu2qfy2exkukuanrkzrrexh9yeg2pr4
  state: repayed
- amount: 1000token
  borrower: cosmos153dk8qh56v4yg6e4uzrvvqjueu6d36fptlr2kw
  collateral: 1000foocoin
  deadline: "20"
  fee: 100token
  id: "1"
  lender: cosmos1qfzpxfhsu2qfy2exkukuanrkzrrexh9yeg2pr4
  state: liquidated
```

```
loand q bank balances $(loand keys show alice -a)
```

`Foocoin` bakiyesi `9000` olarak güncellenmiştir, çünkü `1000foocoin` Alice'in hesabından modül hesabına teminat olarak aktarılmıştır. Alice teminatını kaybetti, ancak kredi tutarını korudu:

```
balances:
- amount: "9000"
  denom: foocoin
- amount: "100000000"
  denom: stake
- amount: "20900"
  denom: token
```

```
loand q bank balances $(loand keys show bob -a)  
```

`Foocoin` bakiyesi `1000` olarak güncellenmiştir, çünkü `1000foocoin` modül hesabından Bob'un hesabına teminat olarak aktarılmıştır. Bob teminatı kazanmıştır, ancak kredi tutarını kaybetmiştir:

```
balances:
- amount: "1000"
  denom: foocoin
- amount: "100000000"
  denom: stake
- amount: "9100"
  denom: token
```
