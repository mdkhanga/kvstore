#ifndef MYSETM_H
#define MYSETM_H

#include <string>

class SetMessage {
	public:
		SetMessage(std::string key, std::string value);
		~SetMessage();
		void serialize(char* buffer);
		static SetMessage deserialize(const char* buffer);
		std::string getKey();
		std::string getValue() ;
	protected:
		std::string key ;
		std::string value ;

} ;

#endif