
#include "spec-testlib.h"

int main(int argc, char *argv[]) {
    setName("sos-checker");
    registerTestlibCmdWithSpecialStream(argc, argv);

    while(!ans.eof()) {
        auto a = ouf.readString();
        auto b = ans.readString();

        if (a != b) {
            quitf(_wa, "wrong answer");
        }
    }

    quitf(_ok, "good");
}

