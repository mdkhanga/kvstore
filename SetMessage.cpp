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

void SetMessage::serialize(char* buffer) {

	std::cout << key.size() << std::endl;
	std::cout << htonl(key.size()) << std::endl; 

	int keysize = key.size();
	memcpy(buffer, &keysize, sizeof(int));
	std::cout << "key size copied to buffer" << std::endl; 
	memcpy(buffer + sizeof(int), key.c_str(), keysize);

	std::cout << "key copied to buffer" << std::endl; 



	int valuesize = value.size();
	memcpy(buffer + sizeof(int) +keysize, &valuesize, sizeof(int));
	memcpy(buffer + sizeof(int) +keysize + sizeof(int) , value.c_str(), valuesize);
	
	std::cout << "value copied to buffer" << std::endl; 

	/* size_t keysize = key.size();
	memcpy(buffer, key.c_str(), keysize);
	size_t valuesize = value.size();
	memcpy(buffer + keysize, value.c_str(), valuesize);

	std::string s2(buffer);
	std::cout << s2 << std::endl;*/


}

SetMessage SetMessage::deserialize(const char* buffer) {
	int strSize;
    memcpy(&strSize, buffer, sizeof(int));
	int ssize = strSize;
	std::cout << ssize << std::endl;
    std::string k(buffer + sizeof(int), ssize);
	std::cout << k << std::endl;

	int vSize;
	memcpy(&vSize, buffer + sizeof(int) + strSize , sizeof(int));
	int vssize = vSize;
	std::cout << vssize << std::endl;
    std::string v(buffer + sizeof(int) + ssize + sizeof(int) , vssize);
	std::cout << v << std::endl;

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