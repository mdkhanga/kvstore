#include<string>

using namespace std;

class SetMessage {
	public:
		SetMessage(string key, string value);
		~SetMessage();
		void serialize(char* buffer);
		static SetMessage deserialize(const char* buffer);
		string getKey();
		string getValue() ;
	protected:
		string key ;
		string value ;

} ;