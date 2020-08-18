package serial

import (
    "context"
    "io"
    "fmt"
    "server/share"
)

func GetFriends(ctx context.Context, user int64) (map[string]*share.User, error) {
    //Produce
    var friendIds []int64
    for it := share.GetFriendIds(user); ; {
        if id, err := it.Next(ctx); err != nil {
            if err == io.EOF {
                break
            }

            return nil, fmt.Errorf("GetFriendIds %d: %w", user, err)
        } else {
            friendIds = append(friendIds, id)
        }
    }

    //Map
    ret := map[string]*share.User{}
    for _, friendId := range friendIds {
        if friend, err := share.GetUserProfile(ctx, friendId); err != nil {
            return nil, fmt.Errorf("GetUserProfile %d: %w", friendId, err)
        } else {
            ret[friend.Name] = friend
        }
    }

    return ret, nil
}
