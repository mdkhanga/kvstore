#ifndef KV_KEYVALUE_H
#define KV_KEYVALUE_H

#include <string>

class KeyValue {
	public:
		KeyValue(std::string key, std::string value);
		~KeyValue();
		int serialize(char* buffer);
		static KeyValue deserialize(const char* buffer);
		std::string getKey();
		std::string getValue() ;
	protected:
		std::string key ;
		std::string value ;
} ;

#endif