// trugw.cpp : This file contains the 'main' function. Program execution begins and ends there.
//

#include <iostream>
#include "..\..\..\..\..\trugw-c\trugw.hpp"

int main()
{
    std::cout << "Hello World!\n";

    std::string socket_path = "/tmp/trugw.sock";
    std::string tru_addr = ":7070";

    // Connect to teogw server
    std::cout << "trying to connect...\n";
    Teogw tgw(socket_path, tru_addr);
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

#include "..\..\..\..\..\trugw-c\trugw.c"