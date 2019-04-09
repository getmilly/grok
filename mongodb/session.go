package mongodb

import (
	"crypto/tls"
	"net"

	mgo "gopkg.in/mgo.v2"
)

//MongoConnection ...
type MongoConnection struct {
	Session *mgo.Session
}

//CreateSession ...
func CreateSession(connectionString string) *MongoConnection {

	session, err := createSession(connectionString)

	if err != nil {
		session, err = createSessionWithSSL(connectionString)

		if err != nil {
			panic(err)
		}
	}

	return &MongoConnection{
		Session: session,
	}
}

func createSession(connectionString string) (*mgo.Session, error) {
	return mgo.Dial(connectionString)
}

func createSessionWithSSL(connectionString string) (*mgo.Session, error) {
	dialInfo, err := mgo.ParseURL(connectionString)

	if err != nil {
		return &mgo.Session{}, err
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		tlsConfig := &tls.Config{}
		return tls.Dial("tcp", addr.String(), tlsConfig)
	}

	return mgo.DialWithInfo(dialInfo)
}
