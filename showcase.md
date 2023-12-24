## Showcase the app

We will showcase the app now (The **api** + **monitoring**)

1. Query the server for resources that have the tag **infra** (nothing comes back)
```bash
curl 'localhost:8000/resources?tags=infra'
```
2. Create 3 resources
    - resource 1 **docker cheat sheet**
        ```bash
        curl -i -XPOST \
        -H 'Content-Type: application/json' \
        -d '{
        "description": "Docker cheat sheet for popular docker commands",
        "reference": "https://github.com/wsargent/docker-cheat-sheet",
        "level": "BEGINNER",
        "type": "ARTICLE",
        "tags": ["docker","infra"] 
         }' localhost:8000/resources
        ```
    - resource 2 **apache kafka documentation**
        ```bash
        curl -i -XPOST \
        -H 'Content-Type: application/json' \
        -d '{
        "description": "Apache kafka article",
        "reference": "https://medium.com/swlh/apache-kafka-what-is-and-how-it-works-e176ab31fcd5",
        "level": "BEGINNER",
        "type": "ARTICLE",
        "tags": ["kafka","messaging","backend"] 
         }' localhost:8000/resources
        ```
    - resource 3 **udemy course of terraform**
        ```bash
        curl -i -XPOST \
        -H 'Content-Type: application/json' \
        -d '{
        "description": "Udemy course showing terraform",
        "reference": "https://www.udemy.com/course/terraform-beginner-to-advanced/",
        "level": "INTERMEDIATE",
        "type": "VIDEO",
        "tags": ["terraform","infra","hashicorp"] 
         }' localhost:8000/resources
        ```
3. Query the service again for **infra** (we get resources 1,3)
    ```bash
    curl 'localhost:8000/resources?tags=infra'
    ```
4. Load test the server using **50,000** concurrent requests from 150 workers.
     Observe the request latencies in grafana.
    ```bash
    bombardier -c 150 -n 10000 'http://localhost:8000/resources?tags=infra'
    ```

5. Now create a very big resource (a lot of text in its description). Load test the app. And notice
    the change in grafana monitoring of the **response size**.
    ```bash
    curl -i -XPOST \
    -H 'Content-Type: application/json' \
    -d '{
    "description": "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industrys standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industrys standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum. Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industrys standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum. Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industrys standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.d web page editors now use Lorem Ipsum as their default model text, and a search for lorem ipsum will uncover many web sites still in their infancy. Various versions have evolved over the years, sometimes by accident,d web page editors now use Lorem Ipsum as their default model text, and a search for lorem ipsum will uncover many web sites still in their infancy. Various versions have evolved over the years, sometimes by accident,d web page editors now use Lorem Ipsum as their default model text, and a search for lorem ipsum will uncover many web sites still in their infancy. Various versions have evolved over the years, sometimes by accident,d web page editors now use Lorem Ipsum as their default model text, and a search for lorem ipsum will uncover many web sites still in their infancy. Various versions have evolved over the years, sometimes by accident,d web page editors now use Lorem Ipsum as their default model text, and a search for lorem ipsum will uncover many web sites still in their infancy. Various versions have evolved over the years, sometimes by accident",
    "reference": "https://medium.com/swlh/apache-kafka-what-is-and-how-it-works-e176ab31fcd5",
    "level": "BEGINNER",
    "type": "ARTICLE",
    "tags": ["kafka","backend","hashicorp"] 
    }' localhost:8000/resources

    ```
