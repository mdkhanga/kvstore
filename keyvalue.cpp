#include "keyvalue.h"
#include <sstream>
#include <iostream>

KeyValue::KeyValue(std::string key, std::string value) {
	this->key = key;
	this->value = value;
}

KeyValue::~KeyValue() {

}

std::string KeyValue::getKey() {
	return key;
}

std::string KeyValue::getValue() {
	return value;
}

int KeyValue::serialize(char* buffer) {

	
	int keysize = key.size();
	int hksize = htonl(keysize);
	memcpy(buffer, &hksize, sizeof(int));
	memcpy(buffer + sizeof(int), key.c_str(), keysize);

	int valuesize = value.size();
	int hvsize = htonl(valuesize);
	memcpy(buffer + sizeof(int) +keysize, &hvsize, sizeof(int));
	memcpy(buffer + sizeof(int) +keysize + sizeof(int) , value.c_str(), valuesize);

	return 4 + keysize + 4 + valuesize ;

}

KeyValue KeyValue::deserialize(const char* buffer) {
	int nstrSize;
    memcpy(&nstrSize, buffer, sizeof(int));
	int ssize = ntohl(nstrSize);
    std::string k(buffer + sizeof(int), ssize);

	int vSize;
	memcpy(&vSize, buffer + sizeof(int) + ssize , sizeof(int));
	int vssize = ntohl(vSize);
    std::string v(buffer + sizeof(int) + ssize + sizeof(int) , vssize);

	return KeyValue(k, v);
}

