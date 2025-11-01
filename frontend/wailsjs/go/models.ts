export namespace config {
	
	export class ActorPhotoConfig {
	    DownloadForKodi: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ActorPhotoConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DownloadForKodi = source["DownloadForKodi"];
	    }
	}
	export class CCConvertConfig {
	    Mode: number;
	    Vars: string;
	
	    static createFrom(source: any = {}) {
	        return new CCConvertConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Mode = source["Mode"];
	        this.Vars = source["Vars"];
	    }
	}
	export class CommonConfig {
	    MainMode: number;
	    SourceFolder: string;
	    FailedOutputFolder: string;
	    SuccessOutputFolder: string;
	    LinkMode: number;
	    ScanHardlink: boolean;
	    FailedMove: boolean;
	    AutoExit: boolean;
	    TranslateToSC: boolean;
	    ActorGender: string;
	    DelEmptyFolder: boolean;
	    NFOSkipDays: number;
	    IgnoreFailedList: boolean;
	    DownloadOnlyMissingImages: boolean;
	    MappingTableValidity: number;
	    Jellyfin: number;
	    ActorOnlyTag: boolean;
	    Sleep: number;
	    AnonymousFill: number;
	    MultiThreading: number;
	    StopCounter: number;
	    RerunDelay: string;
	
	    static createFrom(source: any = {}) {
	        return new CommonConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.MainMode = source["MainMode"];
	        this.SourceFolder = source["SourceFolder"];
	        this.FailedOutputFolder = source["FailedOutputFolder"];
	        this.SuccessOutputFolder = source["SuccessOutputFolder"];
	        this.LinkMode = source["LinkMode"];
	        this.ScanHardlink = source["ScanHardlink"];
	        this.FailedMove = source["FailedMove"];
	        this.AutoExit = source["AutoExit"];
	        this.TranslateToSC = source["TranslateToSC"];
	        this.ActorGender = source["ActorGender"];
	        this.DelEmptyFolder = source["DelEmptyFolder"];
	        this.NFOSkipDays = source["NFOSkipDays"];
	        this.IgnoreFailedList = source["IgnoreFailedList"];
	        this.DownloadOnlyMissingImages = source["DownloadOnlyMissingImages"];
	        this.MappingTableValidity = source["MappingTableValidity"];
	        this.Jellyfin = source["Jellyfin"];
	        this.ActorOnlyTag = source["ActorOnlyTag"];
	        this.Sleep = source["Sleep"];
	        this.AnonymousFill = source["AnonymousFill"];
	        this.MultiThreading = source["MultiThreading"];
	        this.StopCounter = source["StopCounter"];
	        this.RerunDelay = source["RerunDelay"];
	    }
	}
	export class ScraperConfig {
	    Mode: string;
	    MetaTubeURL: string;
	    MetaTubeToken: string;
	    FallbackToLegacy: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ScraperConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Mode = source["Mode"];
	        this.MetaTubeURL = source["MetaTubeURL"];
	        this.MetaTubeToken = source["MetaTubeToken"];
	        this.FallbackToLegacy = source["FallbackToLegacy"];
	    }
	}
	export class STRMConfig {
	    Enable: boolean;
	    PathType: string;
	    ContentMode: string;
	    MultiPartMode: string;
	    NetworkBasePath: string;
	    UseWindowsPath: boolean;
	    ValidateFiles: boolean;
	    StrictValidation: boolean;
	    OutputSuffix: string;
	
	    static createFrom(source: any = {}) {
	        return new STRMConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Enable = source["Enable"];
	        this.PathType = source["PathType"];
	        this.ContentMode = source["ContentMode"];
	        this.MultiPartMode = source["MultiPartMode"];
	        this.NetworkBasePath = source["NetworkBasePath"];
	        this.UseWindowsPath = source["UseWindowsPath"];
	        this.ValidateFiles = source["ValidateFiles"];
	        this.StrictValidation = source["StrictValidation"];
	        this.OutputSuffix = source["OutputSuffix"];
	    }
	}
	export class JellyfinConfig {
	    MultiPartFanart: boolean;
	
	    static createFrom(source: any = {}) {
	        return new JellyfinConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.MultiPartFanart = source["MultiPartFanart"];
	    }
	}
	export class FaceConfig {
	    LocationsModel: string;
	    UncensoredOnly: boolean;
	    AlwaysImagecut: boolean;
	    AspectRatio: number;
	
	    static createFrom(source: any = {}) {
	        return new FaceConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.LocationsModel = source["LocationsModel"];
	        this.UncensoredOnly = source["UncensoredOnly"];
	        this.AlwaysImagecut = source["AlwaysImagecut"];
	        this.AspectRatio = source["AspectRatio"];
	    }
	}
	export class JavdbConfig {
	    Sites: string;
	
	    static createFrom(source: any = {}) {
	        return new JavdbConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Sites = source["Sites"];
	    }
	}
	export class StorylineConfig {
	    Switch: boolean;
	    Site: string;
	    CensoredSite: string;
	    UncensoredSite: string;
	    ShowResult: number;
	    RunMode: number;
	
	    static createFrom(source: any = {}) {
	        return new StorylineConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Switch = source["Switch"];
	        this.Site = source["Site"];
	        this.CensoredSite = source["CensoredSite"];
	        this.UncensoredSite = source["UncensoredSite"];
	        this.ShowResult = source["ShowResult"];
	        this.RunMode = source["RunMode"];
	    }
	}
	export class ExtrafanartConfig {
	    Switch: boolean;
	    ExtrafanartFolder: string;
	    ParallelDownload: number;
	
	    static createFrom(source: any = {}) {
	        return new ExtrafanartConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Switch = source["Switch"];
	        this.ExtrafanartFolder = source["ExtrafanartFolder"];
	        this.ParallelDownload = source["ParallelDownload"];
	    }
	}
	export class WatermarkConfig {
	    Switch: boolean;
	    Water: number;
	
	    static createFrom(source: any = {}) {
	        return new WatermarkConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Switch = source["Switch"];
	        this.Water = source["Water"];
	    }
	}
	export class MediaConfig {
	    MediaType: string;
	    SubType: string;
	
	    static createFrom(source: any = {}) {
	        return new MediaConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.MediaType = source["MediaType"];
	        this.SubType = source["SubType"];
	    }
	}
	export class UncensoredConfig {
	    UncensoredPrefix: string;
	
	    static createFrom(source: any = {}) {
	        return new UncensoredConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.UncensoredPrefix = source["UncensoredPrefix"];
	    }
	}
	export class TrailerConfig {
	    Switch: boolean;
	
	    static createFrom(source: any = {}) {
	        return new TrailerConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Switch = source["Switch"];
	    }
	}
	export class TranslateConfig {
	    Switch: boolean;
	    Engine: string;
	    TargetLang: string;
	    Key: string;
	    Delay: number;
	    Values: string;
	    ServiceSite: string;
	
	    static createFrom(source: any = {}) {
	        return new TranslateConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Switch = source["Switch"];
	        this.Engine = source["Engine"];
	        this.TargetLang = source["TargetLang"];
	        this.Key = source["Key"];
	        this.Delay = source["Delay"];
	        this.Values = source["Values"];
	        this.ServiceSite = source["ServiceSite"];
	    }
	}
	export class DebugModeConfig {
	    Switch: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DebugModeConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Switch = source["Switch"];
	    }
	}
	export class EscapeConfig {
	    Literals: string;
	    Folders: string;
	
	    static createFrom(source: any = {}) {
	        return new EscapeConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Literals = source["Literals"];
	        this.Folders = source["Folders"];
	    }
	}
	export class PriorityConfig {
	    Website: string;
	
	    static createFrom(source: any = {}) {
	        return new PriorityConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Website = source["Website"];
	    }
	}
	export class UpdateConfig {
	    UpdateCheck: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UpdateConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.UpdateCheck = source["UpdateCheck"];
	    }
	}
	export class NameRuleConfig {
	    LocationRule: string;
	    NamingRule: string;
	    MaxTitleLen: number;
	    ImageNamingWithNumber: boolean;
	    NumberUppercase: boolean;
	    NumberRegexs: string;
	
	    static createFrom(source: any = {}) {
	        return new NameRuleConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.LocationRule = source["LocationRule"];
	        this.NamingRule = source["NamingRule"];
	        this.MaxTitleLen = source["MaxTitleLen"];
	        this.ImageNamingWithNumber = source["ImageNamingWithNumber"];
	        this.NumberUppercase = source["NumberUppercase"];
	        this.NumberRegexs = source["NumberRegexs"];
	    }
	}
	export class ProxyConfig {
	    Switch: boolean;
	    Proxy: string;
	    Timeout: number;
	    Retry: number;
	    Type: string;
	    CACertFile: string;
	
	    static createFrom(source: any = {}) {
	        return new ProxyConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Switch = source["Switch"];
	        this.Proxy = source["Proxy"];
	        this.Timeout = source["Timeout"];
	        this.Retry = source["Retry"];
	        this.Type = source["Type"];
	        this.CACertFile = source["CACertFile"];
	    }
	}
	export class Config {
	    Common: CommonConfig;
	    Proxy: ProxyConfig;
	    NameRule: NameRuleConfig;
	    Update: UpdateConfig;
	    Priority: PriorityConfig;
	    Escape: EscapeConfig;
	    DebugMode: DebugModeConfig;
	    Translate: TranslateConfig;
	    Trailer: TrailerConfig;
	    Uncensored: UncensoredConfig;
	    Media: MediaConfig;
	    Watermark: WatermarkConfig;
	    Extrafanart: ExtrafanartConfig;
	    Storyline: StorylineConfig;
	    CCConvert: CCConvertConfig;
	    Javdb: JavdbConfig;
	    Face: FaceConfig;
	    Jellyfin: JellyfinConfig;
	    ActorPhoto: ActorPhotoConfig;
	    STRM: STRMConfig;
	    Scraper: ScraperConfig;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Common = this.convertValues(source["Common"], CommonConfig);
	        this.Proxy = this.convertValues(source["Proxy"], ProxyConfig);
	        this.NameRule = this.convertValues(source["NameRule"], NameRuleConfig);
	        this.Update = this.convertValues(source["Update"], UpdateConfig);
	        this.Priority = this.convertValues(source["Priority"], PriorityConfig);
	        this.Escape = this.convertValues(source["Escape"], EscapeConfig);
	        this.DebugMode = this.convertValues(source["DebugMode"], DebugModeConfig);
	        this.Translate = this.convertValues(source["Translate"], TranslateConfig);
	        this.Trailer = this.convertValues(source["Trailer"], TrailerConfig);
	        this.Uncensored = this.convertValues(source["Uncensored"], UncensoredConfig);
	        this.Media = this.convertValues(source["Media"], MediaConfig);
	        this.Watermark = this.convertValues(source["Watermark"], WatermarkConfig);
	        this.Extrafanart = this.convertValues(source["Extrafanart"], ExtrafanartConfig);
	        this.Storyline = this.convertValues(source["Storyline"], StorylineConfig);
	        this.CCConvert = this.convertValues(source["CCConvert"], CCConvertConfig);
	        this.Javdb = this.convertValues(source["Javdb"], JavdbConfig);
	        this.Face = this.convertValues(source["Face"], FaceConfig);
	        this.Jellyfin = this.convertValues(source["Jellyfin"], JellyfinConfig);
	        this.ActorPhoto = this.convertValues(source["ActorPhoto"], ActorPhotoConfig);
	        this.STRM = this.convertValues(source["STRM"], STRMConfig);
	        this.Scraper = this.convertValues(source["Scraper"], ScraperConfig);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	

}

export namespace gui {
	
	export class FileInfo {
	    path: string;
	    name: string;
	    size: number;
	    number: string;
	    status: string;
	    error: string;
	    duration: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.size = source["size"];
	        this.number = source["number"];
	        this.status = source["status"];
	        this.error = source["error"];
	        this.duration = source["duration"];
	    }
	}
	export class RegexTestRequest {
	    pattern: string;
	    filenames: string[];
	
	    static createFrom(source: any = {}) {
	        return new RegexTestRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pattern = source["pattern"];
	        this.filenames = source["filenames"];
	    }
	}
	export class Stats {
	    total: number;
	    success: number;
	    failed: number;
	    skipped: number;
	    // Go type: time
	    startTime: any;
	    duration: string;
	
	    static createFrom(source: any = {}) {
	        return new Stats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.success = source["success"];
	        this.failed = source["failed"];
	        this.skipped = source["skipped"];
	        this.startTime = this.convertValues(source["startTime"], null);
	        this.duration = source["duration"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace parser {
	
	export class RegexPattern {
	    name: string;
	    pattern: string;
	    description: string;
	    example: string;
	
	    static createFrom(source: any = {}) {
	        return new RegexPattern(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.pattern = source["pattern"];
	        this.description = source["description"];
	        this.example = source["example"];
	    }
	}
	export class RegexTestResult {
	    success: boolean;
	    matched: string;
	    groups: string[];
	    error: string;
	    originalName: string;
	    extractedNumber: string;
	
	    static createFrom(source: any = {}) {
	        return new RegexTestResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.matched = source["matched"];
	        this.groups = source["groups"];
	        this.error = source["error"];
	        this.originalName = source["originalName"];
	        this.extractedNumber = source["extractedNumber"];
	    }
	}

}

