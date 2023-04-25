#include "SetMessage.h"
#include <sstream>
#include <iostream>

SetMessage::SetMessage(std::string key, std::string value) {
	this->key = key;
	this->value = value;
}

SetMessage::~SetMessage() {

}

std::string SetMessage::getKey() {
	return key;
}

std::string SetMessage::getValue() {
	return value;
}

int SetMessage::serialize(char* buffer) {

	
	int keysize = key.size();
	memcpy(buffer, &keysize, sizeof(int));
	memcpy(buffer + sizeof(int), key.c_str(), keysize);

	int valuesize = value.size();
	memcpy(buffer + sizeof(int) +keysize, &valuesize, sizeof(int));
	memcpy(buffer + sizeof(int) +keysize + sizeof(int) , value.c_str(), valuesize);

	return 4 + keysize + 4 + valuesize ;

}

SetMessage SetMessage::deserialize(const char* buffer) {
	int strSize;
    memcpy(&strSize, buffer, sizeof(int));
	int ssize = strSize;
    std::string k(buffer + sizeof(int), ssize);

	int vSize;
	memcpy(&vSize, buffer + sizeof(int) + strSize , sizeof(int));
	int vssize = vSize;
    std::string v(buffer + sizeof(int) + ssize + sizeof(int) , vssize);

	return SetMessage(k, v);
}



/* int main(int argc, char *argv[]) {

	if (argc != 3) {
		cout << "Usage: setm key value" << endl;
		return 0;
	}

	// SetMessage s("key1", "value1") ;
	SetMessage s(argv[1], argv[2]) ;

	char buffer[256] ;

	s.serialize(buffer);

	SetMessage ds = s.deserialize(buffer);

	cout << "key =" << ds.getKey() << endl;
	cout << "value =" << ds.getValue() << endl; 
} */