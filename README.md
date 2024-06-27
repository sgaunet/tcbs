
This golang module offers some basic functions to spin up some containers. 

## Postgresql

```go
  // create the container
  newpgDB, err := tcbs.NewPostgresContainer("postgres", "password", "postgres")
	if err != nil {
		t.Fatalf("could not create postgres container: %v", err)
	}
  // defer the stop of the container
	defer newpgDB.Terminate(context.Background())

  // Open a connection
	db, err := sql.Open("postgres", newpgDB.GetDSNString())
	if err != nil {
		t.Fatalf("could not open postgres connection: %v", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		t.Fatalf("could not ping postgres: %v", err)
	}
```

## redis

```go
  redisC, err := tcbs.NewRedisContainer("", "")
	if err != nil {
		t.Fatalf("could not create redis container: %v", err)
	}
	defer redisC.Terminate(context.Background())

	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Username: redisC.GetRedisUser(),
		Password: redisC.GetRedisPassword(),
		Addr:     redisC.GetRedisHost() + ":" + redisC.GetRedisPort(),
	})
	defer redisClient.Close()

	_, err = redisClient.Ping(ctx).Result()
```