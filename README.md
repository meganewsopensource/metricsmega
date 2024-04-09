# Middleware Metrics

Biblioteca desenvolvida pelo time de Golang da Mega News para melhorar a coleta de informações de requisições.
Essa biblioteca utiliza o Prometheus para coletar as informações.
É utilizada a estrutura de Histogram para personalizar a coleta.
O novo histograma disponibiliza as informações de "code", "method" e "endpoint", além de uma contagem das requisições.
O objetivo dessas métricas é ajudar a visualizar quais requisições estão sendo mais utilizadas e os seus retornos.

Para utilizar tal biblioteca é necessário:
* Importar a biblioteca para o seu projeto.
* Na sua implementação do gin, adicione o middleware com a função Use()

``` golang
r := gin.Default()
r.Use(Config().Metrics)
```

Ao consultar o seu endpoint de metrics, deverá ser possivel encontrar mais uma resposta, o histograma personalizado como o a seguir:

``
"http_request_duration_seconds_bucket{code="200",endpoint="test",method="GET",le="0.1"} 1"
``