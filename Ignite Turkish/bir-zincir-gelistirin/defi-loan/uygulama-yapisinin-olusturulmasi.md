# Uygulama yapısının oluşturulması

Kullanıcıların dijital varlıkları birbirlerine ödünç vermelerini ve birbirlerinden ödünç almalarını sağlayan bir blok zinciri uygulaması için bir yapı oluşturmak üzere, gerekli kodu oluşturmak için Ignite CLI'yı kullanın.

İlk olarak, aşağıdaki komutu çalıştırarak loan adında yeni bir blok zinciri oluşturun:

```
ignite scaffold chain loan --no-module
```

no-module bayrağı Ignite'a varsayılan bir modül oluşturmamasını söyler. Bunun yerine, bir sonraki adımda modülü kendiniz oluşturacaksınız.

Ardından, dizini loan/ olarak değiştirin:

```
cd loan
```

Aşağıdaki komutu çalıştırarak standart Cosmos SDK banka modülüne bağımlı bir modül oluşturun:

```
ignite scaffold module loan --dep bank
```

Özelliklerin listesini içeren bir `loan` modeli oluşturun.

```
ignite scaffold list loan amount fee collateral deadline state borrower lender --no-message
```

no-message bayrağı Ignite'a loan oluşturma, güncelleme ve silme işlemleri için Cosmos SDK mesajları oluşturmamasını söyler. Bunun yerine, özel mesajlar için kod oluşturacaksınız.

Loan'ların talep edilmesi, onaylanması, geri ödenmesi, tasfiye edilmesi ve iptal edilmesine yönelik mesajların işlenmesine yönelik kodu oluşturmak için aşağıdaki komutları çalıştırın:

```
ignite scaffold message request-loan amount fee collateral deadline
```

```
ignite scaffold message approve-loan id:uint
```

```
ignite scaffold message repay-loan id:uint
```

```
ignite scaffold message liquidate-loan id:uint
```

```
ignite scaffold message cancel-loan id:uint
```

Harika bir iş başardınız! Ignite CLI ile birkaç basit komut kullanarak, blok zinciri uygulamanızın temelini başarıyla kurdunuz. Bir loan modeli oluşturdunuz ve mağaza ile etkileşime izin vermek için keeper yöntemlerini dahil ettiniz. Buna ek olarak, beş özel mesaj için mesaj işleyicileri de uyguladınız.

Temel yapı artık yerinde olduğuna göre, geliştirmenin bir sonraki aşamasına geçme zamanı geldi. İlerleyen bölümlerde, oluşturduğunuz mesaj işleyicileri içinde iş mantığını uygulamaya odaklanacaksınız. Bu, her mesaj alındığında gerçekleştirilmesi gereken belirli eylemleri ve işlemleri tanımlamak için kod yazmayı içerecektir.
