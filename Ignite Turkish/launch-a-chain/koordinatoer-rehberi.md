# Koordinatör Rehberi

Koordinatörler Ignite Chain'de yeni zincirler organize eder ve başlatır.

***

Bir zincir başlatma sürecindeki ilk adım, koordinatörün bir zincir başlatma niyetini yayınlamasıdır. `publish` komutu, Ignite üzerinde bir zincir başlatma niyetini bir proje git deposundan yayınlar.

```
ignite n chain publish https://github.com/ignite/example
```

**Çıktı**

```
✔ Source code fetched
✔ Blockchain set up
✔ Chain's binary built
✔ Blockchain initialized
✔ Genesis initialized
✔ Network published
⋆ Launch ID: 3
```

`LaunchID`, Ignite blockchain üzerinde yayınlanan blockchaini tanımlar.

#### Bir genesis belirtin <a href="#specify-a-initial-genesis" id="specify-a-initial-genesis"></a>

Koordinasyon sırasında yeni genesis hesapları ve genesis validatörleri zincir genesis'ine eklenir. Bu hesapların eklendiği ilk oluşum, varsayılan olarak zincir binary'si tarafından oluşturulan varsayılan oluşumdur.

Koordinatör, `--genesis` bayrağı ile zincir başlatma için özel bir başlangıç genesis belirleyebilir. Bu özel ilk oluşum, zincir modülleri için ek varsayılan oluşum hesapları ve özel parametreler içerebilir.

`--genesis-url` bayrağı için bir URL sağlanmalıdır. Bu doğrudan bir JSON genesis dosyasına ya da genesis dosyası içeren bir tarball'a işaret edebilir.

```
ignite n chain publish https://github.com/ignite/example --genesis-url https://raw.githubusercontent.com/ignite/example/master/genesis/gen.json
```

Bir zincirin başlatılması için koordine edilirken, validatörler talep gönderir. Bunlar, zincir için bir validatör olarak oluşumun bir parçası olma taleplerini temsil eder.

Koordinatör bu talepleri listeleyebilir:

> **NOT:** burada "3" `LaunchID`'yi belirtmektedir.

**Çıktı**

```
Id  Status      Type                    Content
1  APPROVED     Add Genesis Account     spn1daefnhnupn85e8vv0yc5epmnkcr5epkqncn2le, 100000000stake
2  APPROVED     Add Genesis Validator   [email protected]:26656, spn1daefnhnupn85e8vv0yc5epmnkcr5epkqncn2le, 95000000stake
3  PENDING      Add Genesis Account     spn1daefnhnupn85e8vv0yc5epmnkcr5epkqncn2le, 95000000stake
4  PENDING      Add Genesis Validator   [email protected]:26656, spn1daefnhnupn85e8vv0yc5epmnkcr5epkqncn2le, 95000000stake
```

Koordinatör bu talepleri onaylayabilir ya da reddedebilir.

Talepleri onaylamak için:

```
ignite n request approve 3 3,4
```

> **NOT:** bir istek listesi seçerken, her iki sentaks da kullanılabilir: `1,2,3,4` ve`1-3,4`.

**Çıktı**

```
✔ Source code fetched
✔ Blockchain set up
✔ Requests format verified
✔ Blockchain initialized
✔ Genesis initialized
✔ Genesis built
✔ The network can be started
✔ Request(s) #3, #4 verified
✔ Request(s) #3, #4 approved
```

Ignite CLI, taleplerin genesis için uygulanabileceğini otomatik olarak doğrular, onaylanan talepler geçersiz bir genesis oluşturmaz.

İstekleri reddetmek için:

```
ignite n request reject 3 3,4
```

**Çıktı**

```
✔ Request(s) #3, #4 rejected
```

***

Oluşum için yeterli sayıda validatör onaylandığında ve koordinatör zincirin başlatılmaya hazır olduğunu düşündüğünde, koordinatör zincirin başlatılmasını başlatabilir.

Bu eylem zincirin oluşumunu sonlandırır, yani zincir için yeni talepler onaylanamaz.

Bu eylem aynı zamanda zincir için başlatma zamanını (veya oluşum zamanını), yani blockchain ağının canlıya geçeceği zamanı da belirler.

**Çıktı**

```
✔ Chain 3 will be launched on 2022-10-01 09:00:00.000000 +0200 CEST
```

Bu örnek çıktı, ağ üzerindeki zincirin başlatılma zamanını gösterir.

#### Özel bir launch zamanı ayarlayın[​](broken-reference) <a href="#set-a-custom-launch-time" id="set-a-custom-launch-time"></a>

Varsayılan olarak, başlatma zamanı mümkün olan en erken tarihe ayarlanacaktır. Uygulamada, doğrulayıcıların düğümlerini ağ lansmanına hazırlamak için zamanları olmalıdır. Bir validatör çevrimiçi olamazsa, validatör setinde hareketsizlik nedeniyle hapse atılabilir.

Koordinatör `--launch-time` bayrağı ile özel bir zaman belirleyebilir.

```
ignite n chain launch --launch-time 2022-01-01T00:00:00Z
```
