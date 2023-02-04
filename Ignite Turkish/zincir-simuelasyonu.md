# Zincir simülasyonu

Ignite CLI zincir simülatörü, mesajları, blokları ve hesapları simüle ederek zinciriniz için bulanıklık testi ve ayrıca kıyaslama testi yapabilmeniz için rastgele girdilere dayalı olarak zincirinizi çalıştırmanıza yardımcı olabilir. Her bir modülde simülasyon testi gerçekleştirmek için bir şablon ve her bir şablonlu mesaj için bir şablon simülasyon yöntemleri oluşturabilirsiniz.

Ignite CLI ile iskele haline getirilen her yeni modül Cosmos SDK [Modül Simülasyonunu](https://docs.cosmos.network/main/building-modules/simulator.html) uygular.

* Her yeni mesaj, testler için gerekli simülasyon yöntemlerini içeren bir dosya oluşturur.
* Scaffolding a `CRUD` type like a `list` or `map` creates a simulation file with `create`, `update`, and `delete` simulation methods in the `x/<module>/simulation` folder and registers these methods in `x/<module>/module_simulation.go`.
* `list` veya `map` gibi bir `CRUD` türünü iskelelemek, `x/<module>/simulation` klasöründe `create`, `update`, ve `delete` simülasyon yöntemleri içeren bir simülasyon dosyası oluşturur ve bu yöntemleri `x/<module>/module_simulation.go` dosyasına kaydeder.
* Tek bir mesajı iskelelemek, kullanıcı tarafından uygulanacak boş bir simülasyon yöntemi oluşturur.

Mesaj tutucu yöntemlerinde yapılan her yeni değişiklik için simülasyon yöntemlerini sürdürmenizi öneririz.

Her simülasyon ağırlıklandırılır çünkü işlemin göndericisi rastgele atanır. Ağırlık, simülasyonun mesajı ne kadar çağıracağını tanımlar.

Daha iyi rastgeleleştirmeler için rastgele bir tohum tanımlayabilirsiniz. Aynı rastgele tohuma sahip simülasyon aynı çıktı ile deterministiktir.

Yeni bir zincir oluşturmak için:

```
ignite scaffold chain mars
```

Bir simülasyonun kayıtlı olmadığını görmek için boş `x/mars/simulation` klasörünü ve `x/mars/module_simulation.go` dosyasını inceleyin.

Şimdi yeni bir mesaj oluşturun:

```
ignite scaffold list user address balance:uint state
```

Yeni bir `x/mars/simulation/user.go` dosyası oluşturulur ve `x/mars/module_simulation.go` dosyasındaki ağırlık ile kaydedilir.

Minimum ağırlık 0 ve maksimum ağırlık 100 olacak şekilde uygun simülasyon ağırlığını tanımladığınızdan emin olun.

Bu örnek için `defaultWeightMsgDeleteUser` değerini 30 ve `defaultWeightMsgUpdateUser` değerini 50 olarak değiştirin.

Tüm modüller için simülasyon testlerini çalıştırmak üzere `app/simulation_test.go` içinde `BenchmarkSimulation` yöntemini çalıştırın:

Simülasyon tarafından sağlanan bayrakları da tanımlayabilirsiniz. Bayraklar `simapp.GetSimulatorFlags()` yöntemi ile tanımlanır:

```
ignite chain simulate -v --numBlocks 200 --blockSize 50 --seed 33
```

Tüm simülasyonun bitmesini bekleyin ve mesajların sonucunu kontrol edin.

Varsayılan `go test` komutu simülasyonu çalıştırmak için çalışır:

```
go test -v -benchmem -run=^$ -bench ^BenchmarkSimulation -cpuprofile cpu.out ./app -Commit=true
```

#### Mesajı atla <a href="#skip-message" id="skip-message"></a>

Hata döndürmeden mesaj göndermekten kaçınmak için mantık kullanın. Simülasyon mesaj işleyicisine yalnızca `simtypes.NoOpMsg(...)` döndürün.

Bir modülü params ile iskelelemek, modülü otomatik olarak `module_simulaton.go` dosyasına ekler:

```
ignite s module earth --params channel:string,minLaunch:uint,maxLaunch:int
```

Parametreler iskele haline getirildikten sonra, `x/<module>/module_simulation.go` dosyasını değiştirerek rastgele parametreleri `RandomizedParams` yöntemine ayarlayın. Simülasyon, fonksiyon çağrısına göre parametreleri rastgele değiştirecektir.

Bir zinciri simüle etmek, [zincir değişmezleri](https://docs.cosmos.network/main/building-modules/invariants.html) hatalarını önlemenize yardımcı olabilir. Değişmez, zincir verilerini geçersiz kılan bir şeyin bozulup bozulmadığını kontrol etmek için zincir tarafından çağrılan bir işlevdir. Yeni bir değişmez oluşturmak ve zincir bütünlüğünü kontrol etmek için, değişmezleri doğrulamak ve tüm değişmezleri kaydetmek üzere bir yöntem oluşturmanız gerekir.

Örneğin, `x/earth/keeper/invariants.go` içinde:

x/earth/keeper/invariants.go

```
package keeper

import (
    "fmt"

    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/tendermint/spn/x/launch/types"
)

const zeroLaunchTimestampRoute = "zero-launch-timestamp"

// RegisterInvariants registers all module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
    ir.RegisterRoute(types.ModuleName, zeroLaunchTimestampRoute,
        ZeroLaunchTimestampInvariant(k))
}

// ZeroLaunchTimestampInvariant invariant that checks if the
// `LaunchTimestamp is zero
func ZeroLaunchTimestampInvariant(k Keeper) sdk.Invariant {
    return func(ctx sdk.Context) (string, bool) {
        all := k.GetAllChain(ctx)
        for _, chain := range all {
            if chain.LaunchTimestamp == 0 {
                return sdk.FormatInvariant(
                    types.ModuleName, zeroLaunchTimestampRoute,
                    "LaunchTimestamp is not set while LaunchTriggered is set",
                ), true
            }
        }
        return "", false
    }
}
```

Şimdi, bekçi değişmezlerini `x/earth/module.go` dosyasına kaydedin:

```
package earth

// ...

// RegisterInvariants registers the capability module's invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
    keeper.RegisterInvariants(ir, am.keeper)
}
```
