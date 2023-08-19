// trugw.cpp : This file contains the 'main' function. Program execution begins and ends there.
//

#include <iostream>
#include "..\..\..\..\..\trugw-c\trugw.hpp"

int main()
{
    // Get tmp path
    size_t buf_count = 256;
    char tmp_path[256];
    getenv_s(&buf_count, tmp_path, "TEMP");
    // char* tmp_path = getenv("TEMP");

    std::string socket_path = std::string(tmp_path) + "\\trugw.sock";
    std::string tru_addr = ":7070";

    std::cout << "Trugw C++ client, sock path: " << socket_path << std::endl;

    // Connect to teogw server
    std::cout << "trying to connect...\n";
    Trugw tgw(socket_path, tru_addr);
    if (!tgw.connected()) {
        std::cout << "can't connect\n";
        return 1;
    }
    std::cout << "connected\n";

    // Send messages
    for (int i = 0; i < 50000; i++) {
        std::string msg = "Hello " + std::to_string(i);
        std::cout << "send " << msg << std::endl;
        tgw.send(msg);

        uint8_t buf[1024];
        auto n = tgw.recv((const char*)buf, sizeof(buf), 0);
        std::string s((const char*)buf, n);
        std::cout << "receive " << s << std::endl;
    }

    return 0;
}

// Run program: Ctrl + F5 or Debug > Start Without Debugging menu
// Debug program: F5 or Debug > Start Debugging menu

// Tips for Getting Started: 
//   1. Use the Solution Explorer window to add/manage files
//   2. Use the Team Explorer window to connect to source control
//   3. Use the Output window to see build output and other messages
//   4. Use the Error List window to view errors
//   5. Go to Project > Add New Item to create new code files, or Project > Add Existing Item to add existing code files to the project
//   6. In the future, to open this project again, go to File > Open > Project and select the .sln file
