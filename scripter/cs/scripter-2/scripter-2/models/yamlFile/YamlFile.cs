using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace scripter.models.yamlFile
{
    public class YamlFile
    {
        public Header Header { get; set; }
        public FileConfiguration FileConfiguration { get; set; }
        public ActionConfiguration ActionConfiguration { get; set; }
        public Context Context { get; set; }
        public Steps Steps { get; set; }
        public Rules Rules { get; set; }
    }
}
