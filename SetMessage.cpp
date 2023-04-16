#include "SetMessage.h"
#include <sstream>
#include <iostream>

SetMessage::SetMessage(string key, string value) {
	this->key = key;
	this->value = value;
}

SetMessage::~SetMessage() {

}

string SetMessage::getKey() {
	return key;
}

string SetMessage::getValue() {
	return value;
}

void SetMessage::serialize(char* buffer) {
	size_t keysize = key.size();
	memcpy(buffer, &keysize, sizeof(size_t));
	buffer = buffer + sizeof(size_t) ;
	memcpy(buffer, key.c_str(), keysize);
	buffer = buffer + keysize ;

	size_t valuesize = value.size();
	memcpy(buffer, &valuesize, sizeof(size_t));
	buffer = buffer + sizeof(size_t) ;
	memcpy(buffer, value.c_str(), valuesize);


}

SetMessage SetMessage::deserialize(const char* buffer) {
	size_t strSize;
    memcpy(&strSize, buffer, sizeof(size_t));
    std::string k(buffer + sizeof(size_t), strSize);

	size_t vSize;
	memcpy(&vSize, buffer + sizeof(size_t) + strSize , sizeof(size_t));
    std::string v(buffer + sizeof(size_t) + strSize + sizeof(size_t) , vSize);

	return SetMessage(k, v);
}



int main() {
	SetMessage s("key1", "value1") ;

	char buffer[40] ;

	s.serialize(buffer);

	SetMessage ds = s.deserialize(buffer);

	cout << "key =" << ds.getKey() << endl;
	cout << "value =" << ds.getValue() << endl;
}