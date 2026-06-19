package config
import(
	"os"
	"flag"
	"log"
	"github.com/ilyakaznacheev/cleanenv"
)
/*
	struct tags: `...`
*/
type HTTPServer struct{
	Addr string `yaml:"address" env-required:"true"`
}

type Config struct{
	Env string `yaml:"env" env:"ENV" env-default:"production"`
	StoragePath string `yaml:"storage-path" env-required:"true"`
	HTTPServer `yaml:"http-server"`
}

func MustLoad() *Config {
	var configPath string
	//To get data from mechine environment veriable
	configPath=os.Getenv("CONFIG_PATH")

	if configPath ==""{
		// flag return pointer, to get data from cmd(argument or flag,using -name) while run the project
		flags:=flag.String("config","","Path to the config file")
		flag.Parse()

		configPath = *flags // dereferencing flag pointer
		if configPath ==""{
			log.Fatal("Config path not set")
		}
	}

	if _,err:=os.Stat(configPath); os.IsNotExist(err){
		log.Fatalf("Config file not exist: %s",configPath)
	}

	var cfg Config
	// Serializing struct config to json,we need to pass pointer(memory addess).
	err:=cleanenv.ReadConfig(configPath,&cfg)

	if err !=nil{
		log.Fatalf("Can not read config file: %s",err.Error())
	}

	return &cfg
}