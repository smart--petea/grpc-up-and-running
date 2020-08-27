#1.grpcurl --import-path proto --proto u.proto -plaintext localhost:10000 proto.ServiceA/homePage
#2.grpcurl --import-path proto --proto u.proto -plaintext localhost:10000 proto.ServiceA/returnAllArticles
#3.grpcurl --import-path proto --proto u.proto -plaintext -d '"1"' localhost:10000 proto.ServiceA/returnSingleArticle
#4.grpcurl --import-path proto --proto u.proto -plaintext -d '{"id":"1"}' localhost:10000 proto.ServiceA/createNewArticle
#5.grpcurl --import-path proto --proto u.proto -plaintext -d '{"id":"1"}' localhost:10000 proto.ServiceA/deleteArticle
