package errgroup

import (
    "server/share"
    "context"
    "io"
    "fmt"
    "sync/atomic"
    "golang.org/x/sync/errgroup"
)

func GetFriends(ctx context.Context, user int64) (map[string]*share.User, error) {
    g, ctx := errgroup.WithContext(ctx)
    friendIds := make(chan int64)

    //Produce
    g.Go(func() error {
        defer close(friendIds)

        for it := share.GetFriendIds(user); ; {
            if id, err := it.Next(ctx); err != nil {
                if err == io.EOF {
                    return nil
                }

                return fmt.Errorf("GetFriendIds %d: %s", user, err)
            } else {
                select {
                case <-ctx.Done():
                    return ctx.Err()
                case friendIds <- id:
                }
            }
        }
    })

    friends := make(chan *share.User)

    nWorkers := 2
    workers := int32(nWorkers)
    for i := 0; i < nWorkers; i++ {
        g.Go(func() error {
            defer func() {
                if atomic.AddInt32(&workers, -1) == 0 {
                    close(friends)
                }
            }()

            for id := range friendIds {
                if friend, err := share.GetUserProfile(ctx, id); err != nil {
                    return fmt.Errorf("GetUserProfile %d: %s", user, err)
                } else {
                    select {
                    case <-ctx.Done():
                        return ctx.Err()
                    case friends <- friend:
                    }
                }
            }

            return nil
        })
    }

    //Reduce
    ret := map[string]*share.User{}
    g.Go(func() error {
        for friend := range friends {
            ret[friend.Name] = friend
        }

        return nil
    })

    return ret, g.Wait()
}
