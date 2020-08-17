package serial

import (
    "context"
    "io"
    "fmt"
    "strconv"
)

type User struct{
    Id int64
    Name string
}

func GetFriends(ctx context.Context, user int64) (map[string]*User, error) {
    //Produce
    var friendIds []int64
    for it := GetFriendIds(user); ; {
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
    ret := map[string]*User{}
    for _, friendId := range friendIds {
        if friend, err := GetUserProfile(ctx, friendId); err != nil {
            return nil, fmt.Errorf("GetUserProfile %d: %w", friendId, err)
        } else {
            ret[friend.Name] = friend
        }
    }

    return ret, nil
}

func GetFriendIds(user int64) *Iterator {
    return NewIterator([]int64 {1, 2, 3, 4})
}

type Iterator struct {
    values []int64
    index int
}

func NewIterator(values []int64) *Iterator {
    return &Iterator{values: values}
}

func (it *Iterator) Next(ctx context.Context) (int64, error) {
    if it.index >= len(it.values) {
        return -1, io.EOF
    }

    value := it.values[it.index]
    it.index = it.index + 1
    return value, nil
}

func GetUserProfile(ctx context.Context, friendId int64) (*User, error) {
    return &User{Name: "friend " + strconv.Itoa(int(friendId))}, nil
}
