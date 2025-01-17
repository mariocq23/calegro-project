using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace scripter.models.yamlFile
{
    internal class Security
    {
        public string User { get; set; }
        public string PublicPassword { get; set; }
        public string PrivatePasswordLocation { get; set; }
        public string CertificateLocation { get; set; }
        public string SecurityDriver { get; set; }
        public bool BypassSecurity { get; set; }
    }
}
