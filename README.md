# GoTodoApp

This is just a simple REST API made in Go which I used to practice writing in Go, and understanding how dependency injection can be done within Go.

## API spec

The API manages "Todo" items defined as below:

```
{
  Id: string 
  Title: string
  Desc: string
  Completed: bool   
}
```

The API supports GET, POST, PUT and DELETE functionality. As this is a simple program validation is kept to a minimum.

As utilizing a real DB service wasn't the purpose of the project the API uses an in-memory DB.
