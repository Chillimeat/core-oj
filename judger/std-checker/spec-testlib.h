

# ifndef SPEC_TESTLIB_H 
# define SPEC_TESTLIB_H

#include "../testlib/testlib.h"

void init_with_special_stream(InStream &ic, char *file_spec, TMode tm) {
    if (!strcmp(file_spec, "stdin")) {
        ic.init(stdin, tm);
    } else if (!strcmp(file_spec, "stdout")) {
        ic.init(stdout, tm);
    } else if (!strcmp(file_spec, "stderr")) {
        ic.init(stderr, tm);
    } else {
        ic.init(file_spec, tm);
    }
}

void registerTestlibCmdWithSpecialStream(int argc, char* argv[]) {
    __testlib_ensuresPreconditions();
    testlibMode = _checker;
    __testlib_set_binary(stdin);

    if (argc > 1 && !strcmp("--help", argv[1]))
        __testlib_help();

    if (argc < 4 || argc > 6)
    {
        quit(_fail, std::string("Program must be run with the following arguments: ") +
            std::string("<input-file> <output-file> <answer-file> [<report-file> [<-appes>]]") + 
            "\nUse \"--help\" to get help information");
    }

    if (argc == 4)
    {
        resultName = "";
        appesMode = false;
    }

    if (argc == 5)
    {
        resultName = argv[4];
        appesMode = false;
    }

    if (argc == 6)
    {
        if (strcmp("-APPES", argv[5]) && strcmp("-appes", argv[5]))
        {
            quit(_fail, std::string("Program must be run with the following arguments: ") +
                        "<input-file> <output-file> <answer-file> [<report-file> [<-appes>]]");
        }
        else
        {
            resultName = argv[4];
            appesMode = true;
        }
    }

    init_with_special_stream(inf, argv[1], _input);
    init_with_special_stream(ouf, argv[2], _output);
    init_with_special_stream(ans, argv[3], _answer);
}

# endif // SPEC_TESTLIB_H