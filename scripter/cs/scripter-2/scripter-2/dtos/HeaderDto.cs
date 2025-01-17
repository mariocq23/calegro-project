using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace scripter_2.dtos
{
    public class HeaderDto
    {
        public string[] import { get; set; }

        public string id { get; set; }

        public string name { get; set; }

        public string inherits { get; set; }

        public string[] implements { get; set; }
    }
}
