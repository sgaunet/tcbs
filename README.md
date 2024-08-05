[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/tcbs)](https://goreportcard.com/report/github.com/sgaunet/tcbs)
[![Maintainability](https://api.codeclimate.com/v1/badges/befd533c3eda78ff851d/maintainability)](https://codeclimate.com/github/sgaunet/tcbs/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/befd533c3eda78ff851d/test_coverage)](https://codeclimate.com/github/sgaunet/tcbs/test_coverage)
[![GoDoc](https://godoc.org/github.com/sgaunet/tcbs?status.svg)](https://godoc.org/github.com/sgaunet/tcbs)
[![License](https://img.shields.io/github/license/sgaunet/tcbs.svg)](LICENSE)

This golang module offers some basic functions to spin up some containers. It uses testcontainers.

## POSTGRESQL

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

## REDIS

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

## SFTP

```go
  container, err := tcbs.NewTestSFTPServer()
	if err!= nil {
		log.Fatal(err)
	}
	defer container.Terminate()
	// init ssh client
	auth := goph.Password(container.GetPassword())
	host := strings.Split(container.GetEndpoint(), ":")[0]
	port := strings.Split(container.GetEndpoint(), ":")[1]
	portUint, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	client, err := goph.NewConn(&goph.Config{
		User:     container.GetUsername(),
		Addr:     host,
		Port:     uint(portUint),
		Auth:     auth,
		Callback: VerifyHost,
	})
	if err != nil {
		log.Fatal(err)
	}
	client.Close()
}

// no fingerprint check
func VerifyHost(host string, remote net.Addr, key ssh.PublicKey) error {
	// return goph.AddKnownHost(host, remote, key, "")
	return nil
}
```

## SMTP

```go
	container, err := tcbs.NewTestSMTPServer()
	if err != nil {
		t.Fatalf("could not create smtp container: %v", err)
	}
	defer container.Terminate()

	// ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	// defer cancel()
	SMTPClient := mail.NewSMTPClient()
	SMTPClient.Host = container.GetSMTPHost()
	SMTPClient.Port = container.GetSMTPPortInt()
	SMTPClient.Username = container.GetSMTPUser()
	SMTPClient.Password = container.GetSMTPPassword()
	// SMTPClient.Encryption = mail.EncryptionTLS
	SMTPClient.Encryption = mail.EncryptionNone
	SMTPClient.KeepAlive = true
	SMTPClient.ConnectTimeout = 10 * time.Second
	SMTPClient.SendTimeout = 10 * time.Second
	SMTPClient.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	client, err := SMTPClient.Connect()
	if err != nil {
		t.Fatalf("Failed to connect to SMTP server: %v", err)
	}
	defer client.Close()
	...
```