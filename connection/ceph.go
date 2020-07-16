package connection

import (
	"fmt"
	"github.com/ceph/go-ceph/rados"
)

type Ceph struct {
	Connection *rados.Conn
	Pools      map[string]bool
}

func NewCeph() (*Ceph, error) {
	conn, err := rados.NewConn()
	if err != nil {
		return nil, err
	}
	ceph := &Ceph{
		Connection: conn,
		Pools:      make(map[string]bool),
	}
	return ceph, nil
}

func (c *Ceph) InitDefault() error {
	err := c.Connection.ReadDefaultConfigFile()
	if err != nil {
		return err
	}
	err = c.Connection.Connect()
	if err != nil {
		return err
	}
	fmt.Println("connect ceph cluster successfully")
	err = c.InitPools()
	if err != nil {
		return err
	}
	return nil
}

const (
	// rgw.bucket.data stores the object
	BucketData = "rgw.bucket.data"

	// rgw.user.uid stores the uid of user and id of user's bucket
	UserUid = "rgw.user.uid"

	// rgw.bucket.index stores the index of bucket
	BucketIndex = "rgw.bucket.index"
)

// InitPools creates the pools the go-rgw needs
func (c *Ceph) InitPools() error {
	existedPools, err := c.Connection.ListPools()
	if err != nil {
		return err
	}
	// determine whether pools already have existed
	for _, value := range existedPools {
		if value == BucketData {
			c.Pools[BucketData] = true
		}
		if value == UserUid {
			c.Pools[UserUid] = true
		}
		if value == BucketIndex {
			c.Pools[BucketIndex] = true
		}
	}
	// if pool doesn't exist, create the pool.
	if _, ok := c.Pools[BucketData]; !ok || !c.Pools[BucketData] {
		err := c.createPool(BucketData)
		if err != nil {
			return err
		}
		c.Pools[BucketData] = true
	}
	if _, ok := c.Pools[UserUid]; !ok || !c.Pools[UserUid] {
		err := c.createPool(UserUid)
		if err != nil {
			return err
		}
		c.Pools[UserUid] = true
	}
	if _, ok := c.Pools[BucketIndex]; !ok || !c.Pools[BucketIndex] {
		err := c.createPool(BucketIndex)
		if err != nil {
			return err
		}
		c.Pools[BucketIndex] = true
	}
	return nil
}

// create pool
func (c *Ceph) createPool(name string) error {
	err := c.Connection.MakePool(name)
	return err
}

// close the connection to the Ceph cluster
func (c *Ceph) Shutdown() {
	if c.Connection != nil {
		c.Connection.Shutdown()
	}
}

func (c *Ceph) WriteObject(pool string, oid string, data []byte, offset uint64) error {
	ioctx, err := c.Connection.OpenIOContext(pool)
	if err != nil {
		return err
	}
	defer ioctx.Destroy()
	err = ioctx.Write(oid, data, offset)
	if err != nil {
		return err
	}
	return nil
}

func (c *Ceph) ReadObject(pool string, oid string, data []byte, offset uint64) (int, error) {
	ioctx, err := c.Connection.OpenIOContext(pool)
	if err != nil {
		return 0, err
	}
	defer ioctx.Destroy()
	num, err := ioctx.Read(oid, data, offset)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (c *Ceph) DeleteObject(pool string, oid string) error {
	ioctx, err := c.Connection.OpenIOContext(pool)
	if err != nil {
		return err
	}
	defer ioctx.Destroy()
	return ioctx.Delete(oid)
}
