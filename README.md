# Cli api maker Tools

## 1.  **Déscription** <br>
   Cli api maker est un outil qui permet de mettre en place rapidement une api complete à partir d'un fichier de configuration <br>
   L'api est faite avec Fastapi et avec une base de donnée sqlite (par defaut)
## 2. **Usage**
### Les commmandes de base : 
- Génerer le fichier de configuration de base : 
  ```
  > cli_api_maker --generate_config
  ```
- Génerer l'api à partir du fichier de config : 
  ```shell
  > cli_api_maler --make_api
  ```
   
### Structure du fichier de config : 
```json
{
    "Api_name": "My api",
    "Api_version": "0.0.0",
    "Api_description": "My api api description",
    "Api_reload": true,
    "Host": "0.0.0.0",
    "Port": "8000",
    "Db": {
        "Db_name": "My api_db",
        "Sgbd": "sqlite",
        "Async": false
    },
    "Models": [
        {
            "Model_name": "todo",
            "Attributs": [
                {
                    "Attribut_name": "name",
                    "Attributs_type": "str"
                }
            ]
        }
    ],
    "Template_dirs": "template", //pas encore necessaire
    "Static_files_dirs": "static_file"//pas encore necessaire mais bientot dispo
}
```
### A propos des configurations : 
- Api_name: le nom de l'api (string)
- Api_version: la version de l'api (selon votre choix)(string)
- Api_description: la description de l'api (true ou false)
- Api_reload: permet d'activer ou non le mode reload auto lors du dev de l'api 
- Host: l'IP de l'api pour le developpement ex 0.0.0.0 ou 127.0.0.1
- Port:  Port occupé par l'api ex 8080 ou 8000
- Db: contient les informations sur la DB qui 
- Db_name: le nom de la DB 
- Sgbd: le sgbd utilisé, choix possible entre mysql, postgresql, ou sqlite(default)
- Models:un tableau qui contient tous les models(class) et peut avoir plusieur éléments 


