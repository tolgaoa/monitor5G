#include <pistache/endpoint.h>
#include <pistache/http.h>
#include <pistache/router.h>

#include "rapidjson/document.h"

#include <iostream>
#include <curl/curl.h>

using namespace Pistache;

void echo(const Rest::Request& req, Http::ResponseWriter resp)
{
    std::string responseString = "Test response\n";
    resp.send(Http::Code::Ok, responseString);
}

static size_t WriteCallback(void *contents, size_t size, size_t nmemb, void *userp)
{
        ((std::string*)userp)->append((char*)contents, size * nmemb);
        return size * nmemb;
}

int main(int argc, char* argv[])
{

	CURL *curl;
        CURLcode res;
        std::string readBuffer;

        struct curl_slist *slist1;
        slist1 = NULL;
        slist1 = curl_slist_append(slist1, "Content-Type: application/json");

        curl = curl_easy_init();

        std::string datasend = "Test message";
	
	std::string uri = "http://20.12.10.12:80/hello";
	
        if(curl) {
                curl_easy_setopt(curl, CURLOPT_URL, uri.c_str());
                curl_easy_setopt(curl, CURLOPT_POST, 1);
                curl_easy_setopt(curl, CURLOPT_HTTPHEADER, slist1);
                curl_easy_setopt(curl, CURLOPT_POSTFIELDS, datasend.c_str());
                //curl_easy_setopt(curl, CURLOPT_POSTFIELDS, json_data);
                //curl_easy_setopt(curl, CURLOPT_READDATA, &readBuffer);
                curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteCallback);
                curl_easy_setopt(curl, CURLOPT_WRITEDATA, &readBuffer);
                res = curl_easy_perform(curl);
                curl_easy_cleanup(curl);
        }

	std::cout << readBuffer << std::endl;

}
