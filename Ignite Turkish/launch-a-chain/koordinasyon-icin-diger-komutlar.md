# Koordinasyon için diğer komutlar

Ignite CLI, zincir başlatmalarını koordine etmek için koordinatörler, validatörler veya diğer katılımcılar tarafından kullanılabilecek çeşitli başka komutlar sunar.

Talepler, validatör katılım talebiyle aynı mantığı izler; oluşumda etkili olabilmeleri için zincir koordinatörü tarafından onaylanmaları gerekir.

***

Herhangi bir katılımcı, zincir için ilişkili bir bakiyeye sahip bir genesis hesabı talep edebilir. Katılımcı, token bakiyelerinin virgülle ayrılmış bir listesini içeren bir adres sağlamalıdır.

Bech32 adresi için herhangi bir önek kullanılabilir, Ignite Chain'de otomatik olarak spn'ye dönüştürülür.

```
ignite n request add-account 3 spn1pe5h2gelhu8aukmrnj0clmec56aspxzuxcy99y 1000stake
```

**Çıktı**

```
Source code fetched
Blockchain set up
⋆ Request 10 to add account to the network has been submitted!
```

***

Herhangi bir katılımcı bir genesis hesabının genesis zincirinden kaldırılmasını talep edebilir. Örneğin, bir kullanıcı ağa zarar verebilecek kadar yüksek bir hesap bakiyesi önerirse bu durum söz konusu olabilir. Katılımcı hesabın adresini vermelidir.

Bech32 adresi için herhangi bir önek kullanılabilir, bu adres Ignite Chain'de otomatik olarak `spn`'ye dönüştürülür.

```
ignite n request remove-account 3 spn1pe5h2gelhu8aukmrnj0clmec56aspxzuxcy99y
```

**Çıktı**

```
Request 11 to remove account from the network has been submitted!
```

***

Herhangi bir katılımcı genesis zincirinden bir genesis validatörünün (gentx) kaldırılmasını talep edebilir. Örneğin, bazı validatörler nedeniyle bir zincir başlatılamadıysa ve bunların genesis'ten kaldırılması gerekiyorsa bu durum söz konusu olabilir. Katılımcı, validatör hesabının adresini sağlamalıdır (genesis hesabıyla aynı formatta).

Bech32 adresi için herhangi bir önek kullanılabilir, Ignite Chain'de otomatik olarak `spn`'ye dönüştürülür.

Talep yalnızca gentx'i genesis'ten kaldırır, ancak ilişkili hesap bakiyesini kaldırmaz.

```
ignite n request remove-validator 429 spn1pe5h2gelhu8aukmrnj0clmec56aspxzuxcy99y
```

**Çıktı**

```
Request 12 to remove validator from the network has been submitted!
```

***
