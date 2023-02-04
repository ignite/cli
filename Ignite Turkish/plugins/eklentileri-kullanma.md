# Eklentileri Kullanma

Ignite eklentileri, Ignite CLI'nin işlevselliğini genişletmek için bir yol sunar. Eklentiler içinde iki temel kavram vardır: `Commands` ve `Hooks`. Commands, cli'nin işlevselliğini genişletirken `Hooks` mevcut komut işlevselliğini genişletir.

Eklentiler, Ignite iskeleli bir Blockchain projesine plugins.yml aracılığıyla veya global olarak `$HOME/.ignite/plugins/plugins.yml` aracılığıyla kaydedilir.

Projenizde bir eklenti kullanmak için proje dizini içinde aşağıdaki komutu çalıştırın:

```
ignite plugin add github.com/project/cli-plugin
```

Eklenti yalnızca proje dizini içinde ignite çalıştırıldığında kullanılabilir olacaktır.

Öte yandan bir eklentiyi global olarak kullanmak için aşağıdaki komutu çalıştırın:

```
ignite plugin add -g github.com/project/cli-plugin
```

Komut, eklentiyi derleyecek ve `ignite`komut listelerinde hemen kullanılabilir hale getirecektir.

Ignite iskeleli bir blockchaindeyken tüm eklentileri ve durumlarını listelemek için `ignite plugin list` komutunu kullanın.

Uzak bir depodaki bir eklenti güncelleme yayınladığında, `ignite plugin update <path/to/plugin>` komutunu çalıştırmak projenizin `config.yml` dosyasında bildirilen belirli bir eklentiyi güncelleyecektir.
