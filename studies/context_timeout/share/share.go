package share

import (
    "context"
    "io"
    "strconv"
)

type User struct{
    Id int64
    Name string
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
