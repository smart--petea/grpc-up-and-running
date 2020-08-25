#1.grpcurl --import-path proto --proto u.proto -plaintext localhost:10000 proto.ServiceA/homePage
#2.grpcurl --import-path proto --proto u.proto -plaintext localhost:10000 proto.ServiceA/returnAllArticles
#3.grpcurl --import-path proto --proto u.proto -plaintext -d '"1"' localhost:10000 proto.ServiceA/returnSingleArticle
