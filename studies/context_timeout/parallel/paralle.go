package parallel

import (
    "server/share"
    "context"
    "io"
    "log"
    "sync/atomic"
)

func GetFriends(ctx context.Context, user int64) (map[string]*share.User, error) {
    friendIds := make(chan int64)

    //Produce
    go func() {
        defer close(friendIds)
        for it := share.GetFriendIds(user); ; {
            if id, err := it.Next(ctx); err != nil {
                if err == io.EOF {
                    break
                }

                //What to do here?
                log.Fatalf("GetFriendIds %d: %s", user, err)
            } else {
                friendIds <- id
            }
        }
    }()

    friends := make(chan *share.User)

    nWorkers := 2
    workers := int32(nWorkers)
    for i := 0; i < nWorkers; i++ {
        go func() {
            defer func() {
                //Last one out closes shop
                if atomic.AddInt32(&workers, -1) == 0 {
                    close(friends)
                }
            }()

            for id := range friendIds {
                if friend, err := share.GetUserProfile(ctx, id); err != nil {
                    //What to do here?
                    log.Fatalf("GetUserProfile %d: %s", user, err)
                } else {
                    friends <- friend
                }
            }
        }()

    }

    //Reduce
    ret := map[string]*share.User{}
    for friend := range friends {
        ret[friend.Name] = friend
    }

    return ret, nil
}
