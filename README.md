# Ozon Tech - Test Task

## Описание проекта

Реализация GraphQL-системы постов и комментариев с возможностью:
- Создания пользователей, постов и комментариев
- Получения постов с сортировкой и пагинацией
- Иерархических запросов получения коммментариев и ответов
- Подписка на создание комментирев к посту

## Что сделано
- Реализовал требуемый функционал
- Разделил приложения на слои (репы, сервисы, gql-резолверы)
- Написал тесты к слою сервисов и реп

## Как сделано
- GraphQL-приложение реализовано с gqlgen(самая популярная и stable-библиотека из того что наресерчил)
- Слои реализуют контракты, зависимости описаны в виде интерфейсов, чтобы было удобнее тестировать
- Тесты написаны с помощью testify+gomock, лежат в одной директории с реализациями
- Паблишер и сабскрайбер реализованы через редис

## Запуск

```bash
# Запуск с PostgreSQL
make run db=pg

# Запуск с in-memory хранилищем
make run
```

## Схема GraphQL
```
scalar Time

type User {
  id: ID!
  username: String!
}

type Post {
  id: ID!
  user: User!
  title: String!
  content: String!
  isCommentable: Boolean!
  comments(limit: Int!, offset: Int!): [Comment!]!
  createdAt: Time!
}

type Comment {
  id: ID!
  user: User!
  parentId: ID
  replies(limit: Int!, offset: Int!): [Comment!]!
  content: String!
  createdAt: Time!
}

enum SortBy {
  NEWEST
  OLDEST
  TOP
}

type Query {
  user(id: ID!): User!
  post(id: ID!): Post!
  replies(commentId: ID!, limit: Int!, offset: Int!): [Comment!]!
  posts(limit: Int!, offset: Int!, sortBy: SortBy): [Post!]!
}

type Mutation {
  createUser(username: String!): User!
  createPost(userId: ID!, title: String!, content: String!, isCommentable: Boolean!): Post!
  createComment(userId: ID!, postId: ID!, parentId: ID, content: String!): Comment!
  togglePostComments(postId: ID!, editor: ID!, enabled: Boolean!): ID!
}

type Subscription {
  commentAdded(postId: ID!): Comment!
}
```

## Примеры запросов
### Создание юзера
```
mutation {
  createUser(username: "ab") {
    id
  }
}
```
### Получение юзера
```
query{
  user(id: "53292ec4-d635-4ef7-a025-8dfd4d485ee1") {
    username
  }
}
```

### Создание поста
```
mutation {
  createPost(
    userId: "53292ec4-d635-4ef7-a025-8dfd4d485ee1"
    title: "a"
    content: "Hi"
    isCommentable: true
  ) {
    id
  }
}
```

### Получение списка постов
```
query{
  posts(limit: 100, offset: 0){
    id
    content
    comments(limit: 2, offset: 0){
      id
      content
      replies(limit: 2, offset: 0){
        id
        content
      }
    }
    user{
      id
    }
  }
}
```
### Получение поста
```
query {
  post(id: "00ccf428-1dc3-4a09-8d75-55be96ba9942") {
    id
    title
    user{
      id
    }
  }
}
```
### Создание комментария
```
mutation {
  createComment(
    userId: "53292ec4-d635-4ef7-a025-8dfd4d485ee1"
    postId: "00ccf428-1dc3-4a09-8d75-55be96ba9942"
    content: "not much"
  ) {
    id
  }
}
```
### Подписка на новые комментарии
```
subscription{
  commentAdded(postId:"00ccf428-1dc3-4a09-8d75-55be96ba9942"){
    id
    content
    parentId
  }
}
```

## Что можно сделать?
- Пересмотреть иерархическую структуру в сторону отдельных запросов для фетча данных
- Добавить полную работу с пермишинами через роли
- Покрыть весь код тестами
